package main

import (
	"fmt"
)

// Image scratch space.
type image struct {
	Xs
	w, h, npx int // height, width, height*width
}

// Clamped value at
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
	for i := 0; i < img.h; i++ {
		s += fmt.Sprintf("%02v: %4.2v\n", i, img.Xs[i*img.w:(i+1)*img.w])
	}
	return s
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
	img := Load("sample.png")

	fmt.Println("Pre:")

	// normalize
	fmt.Println("Normalize:")
	img.stretch()

	// gauss!
	fmt.Println("Blur:")
	Convolve(&img, GaussianKernel(1))

	//
	fmt.Println("Equalize:")
	img.equalize()

	//	// threshold
	//	fmt.Println("Threshold:")
	//	Threshold(img, 0.8)
	//
	//	// count
	//	fmt.Println("Center:")
	//	nx, ny := Count(img)
	//
	//	// blur nx, ny
	//	OneDGauss(nx)
	//	OneDGauss(ny)
	//
	//	// intersect
	//	fmt.Println("Intersect:")
	//	sect := Intersect(nx, ny)
	//	Threshold(sect, 0.8)

	// Export(img, "out.png")
}
