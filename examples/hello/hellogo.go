// Simple hello using DLLCall

//go:generate dllcall helloif.go greeting/greeting.h
package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func main() {
	libname := "greeting.dll"

	err := load_helloif(libname)
	if err != nil {
		log.Fatal(err)
	}
	gr := &greeting{text: fmt.Sprintf("Hello from go %s", time.Now().Format(time.RFC3339)),
		next: &greeting{text: "Next line", next: &greeting{
			text: "Third line",
		}}}
	fillData(gr, 1)
	err = gr.Greet()
	if err != nil {
		log.Fatal(err)
	}
}

func fillData(gr *greeting, idx int) {
	gr.data = []byte(gr.text + " " + strconv.Itoa(idx))
	if gr.next != nil {
		fillData(gr.next, idx+1)
	}
}
