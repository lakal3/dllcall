package main

import (
	"flag"
	"fmt"
	"log"
)

//go:generate dllcall -win msgboxif.go msgboxlib/msgboxif.h

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		usage()
	}
	err := load_msgboxif("msgboxlib.dll")
	if err != nil {
		log.Fatal("Load msgboxlib.dll: ", err)
	}
	mb := msgBox{title: flag.Arg(1), message: flag.Arg(0)}
	err = mb.show()
	if err != nil {
		log.Fatal("Show message: ", err)
	}
}

func usage() {
	fmt.Println("usage: msgbox message title")
}
