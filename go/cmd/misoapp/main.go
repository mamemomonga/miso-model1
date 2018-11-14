package main

import (
	"log"
	"os"
	"time"
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
		hw.LedErr(5)
		time.Sleep(time.Second * 10)
		hw.LedErr(1)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
