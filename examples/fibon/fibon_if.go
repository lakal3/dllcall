package main

/*
#cmethod calc
*/
type calcFibonacci struct {
	n      int64
	result int64
}

/*
#csafe_method fastCalc
*/
type fastcalcFibonacci struct {
	n      int64
	result *int64
}

type extraData struct {
	extras string
}
