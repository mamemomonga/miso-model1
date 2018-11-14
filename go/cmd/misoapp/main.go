package main

import (
	"log"
	"os"
	"time"
	"github.com/mamemomonga/miso-model1/go/hardware"
)

var (
	Version  string
	Revision string
	Basedir  string
)

func main() {
	err := run()
	if err != nil {
		log.Printf("alert: %s", err)
		hw.SLed.Set(hardware.SL_ERR, 5)
		time.Sleep(time.Second * 10)
		hw.SLed.Set(hardware.SL_ERR, 0)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
