// Example to use GLM C++ library to calculate matrix transformation

//go:generate dllcall -fast if.go ../glmcpp/if.h
package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

var count int

func main() {
	flag.IntVar(&count, "count", 1, "Number of iterations")
	debug := flag.Bool("debug", false, "Debug DLL")
	fast := flag.Bool("fast", false, "Debug DLL")
	flag.Parse()
	var err error
	if *debug {
		err = load_if("glmcppd.dll")
	} else {
		err = load_if("glmcpp.dll")
	}
	if err != nil {
		log.Fatal(err)
	}
	mv := MultiplyVectors{}
	mv.Mat[0], mv.Mat[5], mv.Mat[10], mv.Mat[15] = 1, 1, 1, 1
	mv.Mat[12] = 100
	mv.Mat[13] = 50
	for idx := 0; idx < 10; idx++ {
		mv.Vectors = append(mv.Vectors, Vector{float64(idx), float64(idx * 2), float64(idx * 3)})
	}
	st := time.Now()
	for idx := 0; idx < count; idx++ {
		if *fast {
			err = mv.FastMultiply()
		} else {
			err = mv.Multiply()
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	for idx := 0; idx < 5; idx++ {
		fmt.Println(mv.Vectors[idx])
	}
	if count > 1 {
		fmt.Println(" durations ", time.Now().Sub(st).Seconds()*1000, " ms")
	}
}
