package main

import (
	"github.com/IOPaper/Paper/boot"
)

func main() {

	var err error

	if err = boot.Start(); err != nil {
		return
	}

	if err = boot.Shutdown(); err != nil {
		return
	}

}
