package main

import (
	"fmt"
	"os"
)

// Image scratch space.
type image struct {
	Xs
	w, h, npx int // height, width, height*width
}

// Clamped value at

func (img image) Rows() []Xs {
	return img.Chop(img.h)
}

func (img image) at(x, y int) X {
	return img.Xs[clamp(x, img.w)+img.w*clamp(y, img.h)]
}

// Set's value at point (x, y) to v.
func (img image) Set(x, y int, v X) {
	img.Xs[x+img.w*y] = v
}

// Pretty Print
func (img image) String() string {
	s := ""
	rows := img.Rows()
	for i, row := range rows {
		s += fmt.Sprintf("%02v: %v\n", i, row)
	}

	fmt.Println("Stringed")
	return s
}

// Pretty Print
func (xs Xs) String() string {
	return fmt.Sprintf("%4.2v", xs)
}

// Returns a newly allocated image of height h, width w
func New(w, h int) image {
	return image{make(Xs, h*w), w, h, h * w}
}

// SmallTest returns a manually set img for fun.
func SmallTest() image {
	img := New(10, 10)
	img.Xs[15] = 10
	img.Xs[25] = 8
	img.Xs[37] = 12
	img.Set(2, 7, 10)
	img.Set(3, 7, 10)
	img.Set(2, 8, 10)
	img.Set(3, 8, 10)

	return img
}

func main() {
	img := Load(os.Args[1])

	fmt.Println("Pre:")

	// normalize
	fmt.Println("Normalize:")
	//img.stretch()

	// gauss!
	fmt.Println("Blur:")
	//Convolve(&img, LaplaceGaussKernel(1))

	LaplaceGaussKernel(1)

	//
	fmt.Println("Equalize:")
	//img.equalize()

	Export(img, os.Args[2])
}
