package main

import (
	"fmt"
)

func aplay(f string) {
	hw.Aplay(fmt.Sprintf("%s/sounds/%s.wav",Basedir,f))
}
