package main

import (
	"testing"
)

func TestMultiplyVectors_Multiply(t *testing.T) {
	err := load_if("glmcpp.dll")
	if err != nil {
		t.Fatal(err)
		return
	}
	mv := MultiplyVectors{}
	mv.Mat[0], mv.Mat[5], mv.Mat[10], mv.Mat[15] = 1, 1, 1, 1
	mv.Mat[12] = 100
	mv.Mat[13] = 50
	mv.Vectors = append(mv.Vectors, Vector{1, 2, 3})
	err = mv.Multiply()
	if err != nil {
		t.Error(err)
		return
	}
	v := Vector{101, 52, 3}
	if mv.Vectors[0] != v {
		t.Errorf("Assumed 101, 52, 3, not %v", mv.Vectors[0])
	}
}

func TestMultiplyVectors_FastMultiply(t *testing.T) {
	err := load_if("glmcpp.dll")
	if err != nil {
		t.Fatal(err)
		return
	}
	mv := MultiplyVectors{}
	mv.Mat[0], mv.Mat[5], mv.Mat[10], mv.Mat[15] = 1, 1, 1, 1
	mv.Mat[12] = 100
	mv.Mat[13] = 50
	mv.Vectors = append(mv.Vectors, Vector{1, 2, 3})
	err = mv.FastMultiply()
	if err != nil {
		t.Error(err)
		return
	}
	v := Vector{101, 52, 3}
	if mv.Vectors[0] != v {
		t.Errorf("Assumed 101, 52, 3, not %v", mv.Vectors[0])
	}
}

func BenchmarkMultiplyVectors_Multiply(b *testing.B) {
	err := load_if("glmcpp.dll")
	if err != nil {
		b.Fatal(err)
		return
	}
	mv := MultiplyVectors{}
	mv.Mat[0], mv.Mat[5], mv.Mat[10], mv.Mat[15] = 1, 1, 1, 1
	mv.Mat[12] = 100
	mv.Mat[13] = 50
	for n := 0; n < b.N; n++ {
		mv.Vectors = []Vector{Vector{1, 2, 3}}
		err := mv.Multiply()
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkMultiplyVectors_FastMultiply(b *testing.B) {
	err := load_if("glmcpp.dll")
	if err != nil {
		b.Fatal(err)
		return
	}
	mv := MultiplyVectors{}
	mv.Mat[0], mv.Mat[5], mv.Mat[10], mv.Mat[15] = 1, 1, 1, 1
	mv.Mat[12] = 100
	mv.Mat[13] = 50
	for n := 0; n < b.N; n++ {
		mv.Vectors = []Vector{Vector{1, 2, 3}}
		err := mv.FastMultiply()
		if err != nil {
			b.Error(err)
		}
	}
}
