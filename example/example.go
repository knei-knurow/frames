package main

import (
	"fmt"

	"github.com/knei-knurow/frames"
)

func main() {
	// f1 and f2 demostrate an easy to make mistake! mind the difference
	f1 := frames.Create([2]byte{'L', 'D'}, []byte("dondu"))
	fmt.Printf("frame1: %s, len(f1)=%d\n", f1, len(f1))
	for i, v := range f1 {
		fmt.Printf("%d: %s\n", i, frames.DescribeByte(v))
	}

	f2 := []byte("LD5+dondu#q")
	fmt.Printf("frame2: %s, len(f2)=%d\n", f2, len(f2))
	for i, v := range f2 {
		fmt.Printf("%d: %s\n", i, frames.DescribeByte(v))
	}
}
