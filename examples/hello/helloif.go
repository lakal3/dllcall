package main

/*
#cmethod Greet
*/
type greeting struct {
	text string
	data []byte
	next *greeting
}
