package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var fKeep bool
var fFastcall bool
var fLinux bool
var fWindows bool

func main() {
	err := sys_init()
	if err != nil {
		log.Fatal(err)
	}
	flag.BoolVar(&fKeep, "keep", false, "Keep temp generator file")
	flag.BoolVar(&fFastcall, "fast", false, "Use fastcall for csafe_method(s)")
	flag.BoolVar(&fLinux, "linux", false, "Build only Linux interface")
	flag.BoolVar(&fWindows, "win", false, "Build only Windows interface")
	flag.Parse()
	if flag.NArg() < 2 {
		usage()
	}
	src := flag.Arg(0)
	err = parseGoFile(src)
	if err != nil {
		log.Fatal("Parse ", flag.Arg(0), " failed:", err)
	}
	gotarget := src[:len(src)-3] + "_impl"
	err = generate(gotarget, flag.Arg(1))
	if err != nil {
		log.Fatal("Generate ", flag.Arg(1), " failed:", err)
	}
	fmt.Println("Generated: ", flag.Arg(1))
}

func usage() {
	fmt.Println("dllcall v0.12.1")
	fmt.Println("Usage: dllcall [flags] go_file cpp_file")
	os.Exit(1)
}
