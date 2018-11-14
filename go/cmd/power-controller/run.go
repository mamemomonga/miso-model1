package main

import (
	"time"
	"flag"
	"os"
	"os/exec"
//	"log"
)


func run() error {
	action := flag.String("action", "startup", "Action")
	pulse  := flag.Bool("pulse", false, "Startup STATE Pulse")
	flag.Parse()

	switch(*action) {
		case "startup":
			err := run_startup(*pulse)
			if err != nil {
				return err
			}

		case "shutdown":
			err := run_shutdown()
			if err != nil {
				return err
			}

		default:
			flag.PrintDefaults()
			os.Exit(2)
	}
	return nil
}

func run_startup(pulse bool) error {
	if pulse {
		hw.SendStateFlag()
	}
	aplay("silence-1sec")

	hw.SLed.AllOn()
	aplay("dtmf-3");
	time.Sleep(time.Millisecond * 1000)

	for i:=uint8(0); i<16; i++ {
		hw.SLed.Set(i, i*3+2)
	}
	time.Sleep(time.Millisecond * 500)
	aplay("startup");

	hw.SLed.AllOn()
	time.Sleep(time.Millisecond * 1000)
	hw.SLed.AllOff()
	time.Sleep(time.Millisecond * 1000)

	return nil
}


func run_shutdown() error {
	exec.Command("systemctl","stop","misoapp").Run()

	hw.SLed.AllOff()
	time.Sleep(time.Millisecond * 500)
	aplay("shutdown");
	time.Sleep(time.Millisecond * 1000)
	return nil
}

