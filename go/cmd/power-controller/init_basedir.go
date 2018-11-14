package main

import (
	"path/filepath"
	"os"
	"log"
)

func init() {
	b, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	Basedir = b
}

