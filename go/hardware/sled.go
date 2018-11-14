package hardware

// SPI 74HC595
// 16Bit 2連結 LEDアレイ

import (
	"github.com/janne/bcm2835"
	"sync"
	"time"
	"bytes"
//	"log"
)

type SpiLedArray struct {
	led     []uint8     // LED
	counter []uint8     // 点滅用カウンタ
	pin     []uint8     // ピンの状態
	m       *sync.Mutex // Mutex
}

func NewSpiLedArray() (this *SpiLedArray, err error) {
	this = new(SpiLedArray)
	err = nil
	this.led     = make([]uint8, 16)
	this.counter = make([]uint8, 16)
	this.pin     = make([]uint8, 8)
	this.m = new(sync.Mutex)
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


	go func() {
		led := make([]uint8,2)
		for {
			this.m.Lock()
			if bytes.Compare(led, this.led) == 0 {
				this.m.Unlock()
				continue
			}
			led = this.led
			for i:=uint8(0); i<=1; i++ { // 74HC595
				for j:=uint8(0); j<8; j++ { // Pin
					pn := j+8*i
					switch(this.led[pn]) {
						case 0:
							this.pin[i] &=^ ( 1 << j )
						case 1:
							this.pin[i] |= ( 1 << j )
						default:
							if this.counter[pn] == this.led[pn] {
								this.pin[i] ^= ( 1 << j )
								this.counter[pn]=0
							} else {
								this.counter[pn]++
							}
					}
				}
			}
			this.m.Unlock()
			// バイトオーダを逆にする
			bcm2835.SpiTransfern( []byte{ reverse8Bit(this.pin[0]), reverse8Bit(this.pin[1]) } )
			this.m.Unlock()
			time.Sleep(time.Millisecond * 10)
		}
	}()


/*
	pv := make([]uint8,2)
	go func() {
		for {
			this.m.Lock()
			for i:=uint8(0); i<=1; i++ { // 74HC595
				for j:=uint8(0); j<8; j++ { // Pin
					pn := j+8*i
					switch(this.led[pn]) {
						case 0:
							this.pin[i] &=^ ( 1 << j )
						case 1:
							this.pin[i] |= ( 1 << j )
						default:
							if this.counter[pn] == this.led[pn] {
								this.pin[i] ^= ( 1 << j )
								this.counter[pn]=0
							} else {
								this.counter[pn]++
							}
					}
				}
			}
			this.m.Unlock()
			if((pv[0] != this.pin[0]) || (pv[1] != this.pin[1])) {
				this.m.Lock()
				// バイトオーダを逆にする
				bcm2835.SpiTransfern( []byte{ reverse8Bit(this.pin[0]), reverse8Bit(this.pin[1]) } )
				pv[0] = this.pin[0]
				pv[1] = this.pin[1]
				this.m.Unlock()
			}

			time.Sleep(time.Millisecond * 10)
		}
	}()
*/

}

func (this *SpiLedArray) Finalize() {
	bcm2835.SpiEnd()
	bcm2835.Close()
}

func (this *SpiLedArray) Set(ledPin uint8, value uint8) {
	this.m.Lock()
	this.led[ledPin] = value
	this.m.Unlock()
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

