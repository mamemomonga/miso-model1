package main

import (
	"log"
	"time"
	"fmt"
	"strings"
	ha "github.com/mamemomonga/miso-model1/go/hardware"
	"github.com/mamemomonga/misomiso.exe/go/don"
)

const ReportInterval = 15

func run() error {
	hw.SLed.Set(ha.SL_APP, 10)
	log.Printf("info: misolauncher VERSION:%s REVISION:%s", Version, Revision)

	aplay("startup-miso");
	hw.SLed.Set(ha.SL_APP, 1)

	hw.SLed.Set(ha.SL_NET, 5)
	aplay("connecting-mastodon");
	m,err := NewMisoPunch( Basedir+"/config.yaml" )
	if err != nil {
		return err
	}
	defer m.Finish()

	hw.SLed.Set(ha.SL_APP, 1)
	hw.SLed.Set(ha.SL_NET, 1)
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

		hw.SLed.Set(uint8(ha.SL_RUN), 0)
		aplay("finish")
		time.Sleep(time.Millisecond * 100)

		hw.SLed.AllOff()
		hw.SLed.Set(ha.SL_APP, 1)
		hw.SLed.Set(ha.SL_NET, 1)

	}

	return nil
}

func chk_missile_disabled() {
	log.Print("info: CHK_MISSILE_DISABLED")

	if ! hw.Sw_missile_on() {
		return
	}
	hw.SLed.Set(ha.SL_MISSILE, 10)
	hw.SLed.Set(ha.SL_ERR, 1)
	aplay("lockon-already-enable");
	for {
		if ! hw.Sw_missile_on() {
			hw.SLed.Set(ha.SL_MISSILE, 0)
			hw.SLed.Set(ha.SL_ERR,     0)
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
	hw.SLed.Set(ha.SL_MISSILE, 10)

	offmode := false
	for {

		if offmode {
			if hw.Sw_launch_on() {
				hw.SLed.Set(ha.SL_APP, 1)
				hw.SLed.Set(ha.SL_NET, 1)
				hw.SLed.Set(ha.SL_MISSILE, 10)
				hw.SLed.Set(uint8(ha.SL_GAUGE_1 + missile), 10)
				time.Sleep(time.Millisecond * 500)
				offmode = false
			}
		} else {
			if hw.Sw_launch_on() {
				hw.SLed.AllOff()
				offmode = true
				time.Sleep(time.Millisecond * 500)

			}
		}

		if hw.Sw_missile_on() {
			break
		}

		time.Sleep(time.Millisecond * 10)
	}
	hw.Rotary.Stop()

	// せーのっ!
	if hw.Sw_launch_on() && (missile == 6) {
		aplay("seeno");
	}
	log.Printf("MISSILE: %d",missile)

	hw.SLed.Set(uint8(ha.SL_GAUGE_1 + missile), 1)
	time.Sleep(time.Millisecond * 100)

	aplay("missile-type");
	aplay(""+target_sounds[missile]);
	aplay("is-armed");

}

func lockon() bool {
	hw.SLed.Set(ha.SL_MISSILE, 5)
	log.Print("info: lockon")

	if ! hw.Sw_missile_on() {
		aplay("lockon-disabled");
		hw.SLed.Set(uint8(ha.SL_GAUGE_1 + missile), 0)
		hw.SLed.Set(ha.SL_READY, 0)
		return false
	}
	hw.SLed.Set(ha.SL_MISSILE, 1)
	aplay("lockon");
	hw.SLed.Set(ha.SL_READY, 10)

	for {
		if ! hw.Sw_missile_on() {
			aplay("lockon-disabled");
			hw.SLed.Set(uint8(ha.SL_GAUGE_1 + missile), 0)
			hw.SLed.Set(ha.SL_READY, 0)
			return false
		}
		if hw.Sw_launch_on() {
			hw.SLed.Set(uint8(ha.SL_READY), 2)
			hw.SLed.Set(uint8(ha.SL_RUN), 2)
			aplay("launch")
			hw.SLed.Set(uint8(ha.SL_READY), 1)
			hw.SLed.Set(uint8(ha.SL_RUN), 20)
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
			if hw.Sw_launch_on() {
				this.d.Toot(fmt.Sprintf("%s #misomiso", this.Keyword))
				hw.SLed.Set(uint8(ha.SL_READY), 10)
				aplay(target_sounds[missile])
				hw.SLed.Set(uint8(ha.SL_READY), 1)
			}
			// ミサイルスイッチオフで中断
			if ! hw.Sw_missile_on() {
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
	hw.SLed.Set(uint8(ha.SL_RUN), 2)
	aplay(target_sounds[missile])
	aplay("hit")
	time.Sleep(time.Millisecond * 100)
	hw.SLed.Set(uint8(ha.SL_RUN), 20)
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

