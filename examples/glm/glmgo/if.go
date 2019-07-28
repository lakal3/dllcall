//

package main

type Vector [3]float64

type Matrix4x4 [16]float64

/*
#cmethod Multiply
#csafe_method FastMultiply
*/
type MultiplyVectors struct {
	Mat     Matrix4x4
	Vectors []Vector
}
