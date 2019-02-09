// Simple hello using DLLCall

//go:generate dllcall helloif.go greeting/greeting.h
package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	err := load_helloif("greeting.dll")
	if err != nil {
		log.Fatal(err)
	}
	gr := &greeting{greeting: fmt.Sprintf("Hello DLLCall!. Time is %v", time.Now())}
	err = gr.Greet()
	if err != nil {
		log.Fatal(err)
	}
}
