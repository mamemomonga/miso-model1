package hardware

// MISO MODEL1

import (
	"time"
	"os/exec"
	"github.com/stianeikeland/go-rpio"
//	"log"
)

// SPI LED
const (
	SL_MISSILE = 1
	SL_RUN     = 3
	SL_READY   = 4
	SL_ERR     = 5
	SL_NET     = 6
	SL_APP     = 7
	SL_GAUGE_1 = 8
	SL_GAUGE_2 = 9
	SL_GAUGE_3 = 10
	SL_GAUGE_4 = 11
	SL_GAUGE_5 = 12
	SL_GAUGE_6 = 13
	SL_GAUGE_7 = 14
	SL_GAUGE_8 = 15
)

// GPIO
const (
	G_SPEAKER_AMP = 12
	G_MISSILE     = 5
	G_LAUNCH      = 13
	G_ROTARY_A    = 26
	G_ROTARY_B    = 19
	G_STATE       = 16
)

type Hardware struct {
	SLed        *SpiLedArray
	Rotary      *EC12
	PinMissile  rpio.Pin
	PinLaunch   rpio.Pin
}

func NewHardware()(this *Hardware, err error) {
	this = new(Hardware)

	err = rpio.Open()
	if err != nil {
		return
	}

	this.SLed,err = NewSpiLedArray()
	if err != nil {
		return
	}

	this.Rotary, err = NewEC12(G_ROTARY_A, G_ROTARY_B)
	if err != nil {
		return
	}

	this.PinMissile = rpio.Pin(G_MISSILE)
	this.PinMissile.PullDown()
	this.PinMissile.Input()

	this.PinLaunch = rpio.Pin(G_LAUNCH)
	this.PinLaunch.PullUp()
	this.PinLaunch.Input()

	this.SLed.Run()
	this.SLed.AllOff()
	return
}

func (this *Hardware) Aplay(filename string)(err error) {
	pin := rpio.Pin(G_SPEAKER_AMP)
	pin.Output()
	pin.High()
	time.Sleep(time.Millisecond * 10)

	err = exec.Command("aplay", filename).Run()
	if err != nil {
		return
	}

	time.Sleep(time.Millisecond * 10)
	pin.Low()
	return
}

func (this *Hardware) SendStateFlag() {
	pin := rpio.Pin(G_STATE)
	pin.Output()
	pin.Low()
	time.Sleep(time.Millisecond * 100)
	pin.Input() // HiZ
}

func (this *Hardware) Sw_missile_on()(ret bool) {
	// 正論理
	if this.PinMissile.Read() == 0 {
		return false
	} else {
		return true
	}
}

func (this *Hardware) Sw_launch_on()(ret bool) {
	// 負論理
	if this.PinLaunch.Read() == 0 {
		return true
	} else {
		return false
	}
}

func (this *Hardware) RotarySelector(retval *int) {
	this.Rotary.Range( EC12Range{
		Start: 0, Max: 7, Min: 0,
		Selected: func(val int) {
			this.SLed.Set(uint8(val+SL_GAUGE_1),10)
			*retval = val
		},
		Clear: func(val int) {
			for i:= uint8(0); i<8; i++ {
				this.SLed.Set(uint8(i+SL_GAUGE_1),0)
			}
		},
	})
	return
}

func (this *Hardware) Finalize() {
	this.SLed.AllOff()
	rpio.Close()
	return
}


