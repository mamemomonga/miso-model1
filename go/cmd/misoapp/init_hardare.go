package main

import (
	"log"
	"github.com/mamemomonga/miso-model1/go/hardware"
)

var hw *hardware.Hardware

func init() {
	h, err := hardware.NewHardware()
	if err != nil {
		log.Fatal(err)
	}
	hw = h
}
