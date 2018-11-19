package hardware

// SPI 74HC595
// 16Bit 2連結 LEDアレイ

import (
	"github.com/janne/bcm2835"
	"sync"
	"time"
//	"log"
)

type SpiLedArray struct {
	led           []uint8     // LED(0: 消灯 1:点灯 >1:点滅)
	counter       []uint8     // 点滅用カウンタ
	pin           []uint8     // ピンの状態
	m             *sync.Mutex // Mutex
	sleep_timer   uint64
}

func NewSpiLedArray() (this *SpiLedArray, err error) {
	this = new(SpiLedArray)
	err = nil
	this.led     = make([]uint8, 16)
	this.counter = make([]uint8, 16)
	this.pin     = make([]uint8, 2)
	this.m       = new(sync.Mutex)
	return
}

func (this *SpiLedArray) Run() {
	bcm2835.Init()
	bcm2835.SpiBegin() // SPI0
	bcm2835.SpiSetBitOrder(BCM2835_SPI_BIT_ORDER_MSBFIRST) // 効かない？
	bcm2835.SpiSetDataMode(BCM2835_SPI_MODE2)
	bcm2835.SpiSetClockDivider(BCM2835_SPI_CLOCK_DIVIDER_128)
	bcm2835.SpiChipSelect(BCM2835_SPI_CS0)
	bcm2835.SpiSetChipSelectPolarity(BCM2835_SPI_CS0, LOW)

	// 点滅
	go this.blinker()
}


func (this *SpiLedArray) blinker() {
	for {
		flag := false
		this.m.Lock()
		for i:=uint8(0); i<16; i++ {
			if this.counter[i] == 0 {
				continue
			}
			flag = true
			if this.led[i] == this.counter[i] {
				c,p := this.l2p(i)
				this.pin[c] ^= ( 1 << p )
				this.counter[i]=1
			}
			this.counter[i]++
		}
		if flag {
			this.update()
		}
		this.m.Unlock()
		time.Sleep(time.Millisecond * 10)
	}
}

func (this *SpiLedArray) l2p(led uint8)(chip uint8, value uint8) {
	if led < 8 {
		chip  = 0
		value = led
	} else {
		chip   = 1
		value = led-8
	}
	return chip,value
}


func (this *SpiLedArray) update() {
	bcm2835.SpiTransfern( []byte{ reverse8Bit(this.pin[0]), reverse8Bit(this.pin[1]) } )
}

func (this *SpiLedArray) Finalize() {
	bcm2835.SpiEnd()
	bcm2835.Close()
}

func (this *SpiLedArray) Set(led uint8, val uint8) {
	c,p := this.l2p(led)
	this.m.Lock()
	this.led[led]=val
	switch(val) {
		case 0:
			this.pin[c] &=^ ( 1 << p )
			this.counter[led]=0
			this.update()
		case 1:
			this.pin[c] |=  ( 1 << p )
			this.counter[led]=0
			this.update()
		default:
			this.pin[c] |=  ( 1 << p )
			this.counter[led]=val
	}
	this.m.Unlock()
}

func (this *SpiLedArray) SleepTimer(s int) {
}

func (this *SpiLedArray) AllOff() {
	for i:=uint8(0); i<16; i++ {
		this.Set(i,0)
	}
}

func (this *SpiLedArray) AllOn() {
	for i:=uint8(0); i<16; i++ {
		this.Set(i,1)
	}
}

