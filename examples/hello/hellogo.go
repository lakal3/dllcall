// Simple hello using DLLCall

//go:generate dllcall helloif.go greeting/greeting.h
package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	libname := "greeting.dll"

	err := load_helloif(libname)
	if err != nil {
		log.Fatal(err)
	}
	gr := &greeting{greeting: fmt.Sprintf("Hello DLLCall!. Time is %v", time.Now())}
	err = gr.Greet()
	if err != nil {
		log.Fatal(err)
	}
}
