package hardware

// ロータリーエンコーダ
// ALPS EC12E2420801
// http://akizukidenshi.com/catalog/g/gP-06357/

import (
	"time"
	"sync"
	"github.com/stianeikeland/go-rpio"
//	"log"
)

type EC12Range struct {
	Start int
	Max   int
	Min   int
	Selected func(int)
	Clear  func(int)
}

type EC12 struct {
	pinA uint8
	pinB uint8
	stop bool
	m *sync.Mutex
}

func NewEC12(pA uint8, pB uint8) (this *EC12, err error) {
	this = new(EC12)
	this.m = new(sync.Mutex)
	this.pinA = pA
	this.pinB = pB
	this.stop = false
	err = nil
	return
}

func (this *EC12) Selector(fInc func(), fDec func())(err error) {
	pinA := rpio.Pin(this.pinA)
	pinB := rpio.Pin(this.pinB)
	pinA.Input()
	pinA.PullUp()
	pinA.Input()
	pinB.PullUp()

	//  Bits: XXXXCDAB | A,B: 現在値 C,D: 前回値
	rc := uint8(0)

	go func() {
		for {
			this.m.Lock()
			if this.stop {
				this.stop = false
				break
			}
			this.m.Unlock()
			stA := pinA.Read()
			stB := pinB.Read()
			rc = ( rc << 2 ) + uint8((stA << 1) + stB) // 前回の値をシフトして今回の値を加算
			rc &= 15 // ビットマスク 00001111
			if (rc & 3) != (rc >> 2){ // 前回と今回の値が違う
				// log.Printf("A:%d B:%d RC:%02d",stA, stB, rc)
				if rc == 2 {
					fInc()
				} else if rc == 7 {
					fDec()
				}
			}
			time.Sleep(time.Millisecond * 5)
		}
		this.m.Unlock()
	}()
	return err
}

func (this *EC12) Stop() {
	this.m.Lock()
	this.stop = true
	this.m.Unlock()
}

func (this *EC12) Range(rs EC12Range)(err error) {
	val := rs.Start
	rs.Selected(val)
	return this.Selector(
		// 増加
		func() {
			rs.Clear(val)
			if val == rs.Max {
				val = rs.Min
			} else {
				val++
			}
			// log.Print("INC",val)
			rs.Selected(val)
		},
		// 減少
		func() {
			rs.Clear(val)
			if val == 0 {
				val = rs.Max
			} else {
				val--
			}
			// log.Print("DEC",val)
			rs.Selected(val)
		},
	)
}
