// Fibonacci example to measure performance between different calls

//go:generate dllcall -fast fibon_if.go fibonlib/fibon_if.h
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var count int

func main() {
	flag.IntVar(&count, "count", 1, "Number of iterations")
	flag.Parse()
	if flag.NArg() < 2 {
		usage()
	}
	value, err := strconv.ParseInt(flag.Arg(1), 10, 64)
	if err != nil {
		log.Fatal("Invalid value: ", err)
	}
	var result *int64
	start := time.Now()
	dllName := "fibonlib.dll"
	err = load_fibon_if(dllName)
	if err != nil {
		log.Fatal("Failed to load ", dllName, ": ", err)
	}
	var fn func(int64) (*int64, error)
	switch strings.ToLower(flag.Arg(0)) {
	case "go":
		fn = goFibon
	case "syscall":
		fn = stdFibon
	case "fastcall":
		fn = fastFibon
	default:
		usage()
	}
	for idx := 0; idx < count; idx++ {
		result, err = fn(value)
		if err != nil {
			log.Fatal("Calculation failure: ", err)
		}
	}
	fmt.Print("Result: ", *result)
	if count > 1 {
		fmt.Println(" durations ", time.Now().Sub(start).Seconds()*1000, " ms")
	} else {
		fmt.Println()
	}
}

func goFibon(n int64) (*int64, error) {
	if n < 1 {
		return nil, errors.New("Value must be at least 1")
	}
	r := fibon(n)
	return &r, nil
}

func stdFibon(n int64) (*int64, error) {
	var res int64
	cl := &calcFibonacci{n: n, result: &res}
	err := cl.calc()
	if err != nil {
		return nil, err
	}
	return cl.result, nil
}

func fastFibon(n int64) (*int64, error) {
	cl := &calcFibonacci{n: n, result: new(int64)}
	err := cl.fastCalc()
	if err != nil {
		return nil, err
	}
	return cl.result, nil
}

func fibon(n int64) int64 {
	if n > 2 {
		return fibon(n-1) + fibon(n-2)
	}
	return 1
}

func usage() {
	fmt.Println("fibon -count n method value")
	fmt.Println("  value - Fibonacci to calculate")
	fmt.Println("  count - Number of iterations (is set, total runtime will be shown)")
	fmt.Println("  methods: go - go fibonacci")
	fmt.Println("           syscall - c fibonacci using standard call")
	fmt.Println("           fastcall - c fibonacci using fast call")
	os.Exit(1)
}
