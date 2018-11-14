package main

import (
	"log"
	"time"
	"fmt"
	"strings"
	"github.com/mamemomonga/misomiso.exe/go/don"
)

const ReportInterval = 15

func run() error {
	hw.LedApp(10)
	log.Printf("info: misolauncher VERSION:%s REVISION:%s", Version, Revision)

	aplay("startup-miso");
	hw.LedApp(1)

	hw.LedNet(5)
	aplay("connecting-mastodon");
	m,err := NewMisoPunch( Basedir+"/config.yaml" )
	if err != nil {
		return err
	}
	defer m.Finish()

	hw.LedApp(1)
	hw.LedNet(1)
	aplay("connected");

	// ------------------------------

	for {
		chk_missile_disabled()
		select_missile()
		if b := lockon(); b==false {
			continue
		}
		m.Keyword = target_text[missile]
		m.Regexp  = target_regexp[missile]

		err = m.Run()
		if err != nil {
			return err
		}

		hw.LedRun(0)
		aplay("finish")
		time.Sleep(time.Millisecond * 100)

		hw.LedAllOff()
		hw.LedApp(1)
		hw.LedApp(1)
	}

	return nil
}

func chk_missile_disabled() {
	log.Print("info: CHK_MISSILE_DISABLED")

	if ! hw.SwMissileOn() {
		return
	}
	hw.LedMissile(10)
	hw.LedErr(1)
	aplay("lockon-already-enable");
	for {
		if ! hw.SwMissileOn() {
			hw.LedMissile(0)
			hw.LedErr(0)
			aplay("lockon-disabled");
			return
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func select_missile() {
	log.Print("info: SELECT_MISSILE")
	missile = 0

	hw.RotarySelector(&missile)
	aplay("select-missile");
	hw.LedMissile(10)

	offmode := false
	for {

		if offmode {
			if hw.SwLaunchOn() {
				hw.LedApp(1)
				hw.LedNet(1)
				hw.LedMissile(10)
				hw.LedGauge(missile, 10)
				time.Sleep(time.Millisecond * 500)
				offmode = false
			}
		} else {
			if hw.SwLaunchOn() {
				hw.LedAllOff()
				offmode = true
				time.Sleep(time.Millisecond * 500)

			}
		}

		if hw.SwMissileOn() {
			break
		}

		time.Sleep(time.Millisecond * 10)
	}
	hw.Rotary.Stop()

	// せーのっ!
	if hw.SwLaunchOn() && (missile == 6) {
		aplay("seeno");
	}
	log.Printf("MISSILE: %d",missile)

	hw.LedGauge(missile, 1)
	time.Sleep(time.Millisecond * 100)

	aplay("missile-type");
	aplay(""+target_sounds[missile]);
	aplay("is-armed");

}

func lockon() bool {
	hw.LedMissile(5)
	log.Print("info: lockon")

	if ! hw.SwMissileOn() {
		aplay("lockon-disabled");
		hw.LedGauge(missile, 0)
		hw.LedReady(0)
		return false
	}
	hw.LedMissile(1)
	aplay("lockon");
	hw.LedReady(10)

	for {
		if ! hw.SwMissileOn() {
			aplay("lockon-disabled");
			hw.LedGauge(missile, 0)
			hw.LedReady(0)
			return false
		}
		if hw.SwLaunchOn() {
			hw.LedReady(2)
			hw.LedRun(2)
			aplay("launch")
			hw.LedReady(1)
			hw.LedRun(20)
			return true
		}
		time.Sleep(time.Millisecond * 10)
	}
	return false
}

// -----------------------

type MisoPunch struct {
	d       *don.Don
	l       *don.Launcher
	Keyword string
	Regexp  string
}

func NewMisoPunch(config_file string)(this *MisoPunch,err error) {
	this = new(MisoPunch)

	d, err := don.NewDon(config_file)
	if err != nil {
		return this,err
	}
	this.d = d

	err = this.d.Connect()
	if err != nil {
		return this,err
	}
	return this, nil
}

func (this *MisoPunch) Run() error {
	l, err := this.d.Launcher( don.LauncherConf{
		SearchRegexp: this.Regexp,
		Callbacks : don.LauncherCallbacks{
			Boost:   this.cBoosted,
			Timeout: this.cReportFinish,
			Abort:   this.cReportFinish,
		},
	})
	if err != nil {
		return err
	}
	this.l = l

	err = this.d.Toot(fmt.Sprintf("[発射]\n%s #misomiso", this.Keyword))
	if err != nil {
		return err
	}

	this.l.Run()

	// 中間報告
	go func() {
		time.Sleep(time.Second * ReportInterval)
		for {
			if ! this.l.IsRunning() {
				break
			}
			this.cReportRunning( this.l.Report() )
			time.Sleep(time.Second * ReportInterval)
		}
	}()

	go func() {
		for {
			if ! this.l.IsRunning() {
				break
			}
			// 発射ボタンで追加トゥート
			if hw.SwLaunchOn() {
				this.d.Toot(fmt.Sprintf("%s #misomiso", this.Keyword))
				hw.LedReady(10)
				aplay(target_sounds[missile])
				hw.LedReady(1)
			}
			// ミサイルスイッチオフで中断
			if ! hw.SwMissileOn() {
				this.l.Abort()
			}
			time.Sleep(time.Millisecond * 10)
		}
	}()

	err = this.l.Reciever()
	if err != nil {
		return err
	}

	return nil
}

func (this *MisoPunch) Finish() {
	this.l.Finish()
}

func (this *MisoPunch) cBoosted() {
	log.Print("info: 命中")
	hw.LedRun(2)
	aplay(target_sounds[missile])
	aplay("hit")
	time.Sleep(time.Millisecond * 100)
	hw.LedRun(20)
}

func (this *MisoPunch) cReportFinish(r don.LauncherReport) {

	mbs := []string{}
	for _,i := range r.Members {
		mbs = append(mbs, i+" さん")
	}
	mb := strings.Join( mbs,", ")

	e := fmt.Sprintf("%.0f分%02d秒",r.Elapsed.Truncate(time.Minute).Minutes(), int(r.Elapsed.Seconds()) % 60 )

	var n string
	if r.Hit == 0 {
		n = fmt.Sprintf("[自爆]\n戦果: %d%s\n飛行時間: %s\n追尾終了しました。#misomiso",
		r.Hit, this.Keyword, e )

	} else {
		n = fmt.Sprintf("[自爆]\n戦果: %d%s\n飛行時間: %s\n追尾終了しました。\n%s ありがとうございました。#misomiso",
			r.Hit, this.Keyword, e, mb )
	}

	this.d.Toot(n)
}

func (this *MisoPunch) cReportRunning(r don.LauncherReport) {
	m := fmt.Sprintf("%.0f分%02d秒",r.Remain.Truncate(time.Minute).Minutes(),  int(r.Remain.Seconds()) % 60 )
	e := fmt.Sprintf("%.0f分%02d秒",r.Elapsed.Truncate(time.Minute).Minutes(), int(r.Elapsed.Seconds()) % 60 )
	n := fmt.Sprintf("[追尾中]\n戦果: %d%s\n残り時間: %s\n経過時間: %s\n#misomiso", r.Hit, this.Keyword, m, e)
	this.d.Toot(n)
}

