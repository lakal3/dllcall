package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var fKeep bool

func main() {
	err := sys_init()
	if err != nil {
		log.Fatal(err)
	}
	flag.BoolVar(&fKeep, "keep", false, "Keep temp generator file")
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
	fmt.Println("dllcall v0.7.1")
	fmt.Println("Usage: dllcall [flags] go_file cpp_file")
	os.Exit(1)
}
