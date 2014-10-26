package main

import (
	"fmt"
	"math"
	"runtime"
)

type X float64
type Xs []X

// Implements Stringer.
// Pretty Print: hell, because though Xs lets me do cool things,
// it also blocks any builting fmting
func (xs Xs) String() (s string) {
	for _, v := range xs {
		s += fmt.Sprintf("|%5.2v ", v)
	}
	return s + "|"
}

func (xs Xs) StringChop(n int) (s string) {
	rows := xs.Chop(n)
	for i, row := range rows {
		s += fmt.Sprintf("%02v: %v\n", i, row)
	}
	return s
}

// implements color.Color
func (s X) RGBA() (r, g, b, a uint32) {
	x := uint32(s * 0xFFFF)
	return x, x, x, 0xFFFF
}

// n is how many chunks you want
func (xs Xs) Chop(n int) []Xs {
	xss := make([]Xs, len(xs)/n)
	for i := range xss {
		xss[i] = xs[i*n : (i+1)*n]
	}
	return xss
}

func (xs Xs) Split(N int, f func(Xs, int)) {
	n := len(xs)
	done := make(chan struct{})

	for p := 0; p < N; p++ {
		go func(p int) {
			offset := p * n / N
			slice := xs[offset : (p+1)*n/N]
			f(slice, offset)
			done <- struct{}{}
		}(p)
	}

	for c := 0; c < N; c++ {
		<-done
	}

	close(done)
	return
}

func (r Xs) extrema() (min, max X) {
	min, max = math.MaxFloat64, -math.MaxFloat64
	for _, x := range r {
		if x < min {
			min = x
		}
		if x > max {
			max = x
		}
	}
	return
}

// Normalize shifts by the minimum, and scales everything to [0, 1.0]
func (xs Xs) stretch() {
	min, max := xs.extrema()
	shift, scale := min, 1/(max-min)

	xs.Split(runtime.NumCPU(), func(slice Xs, _ int) {
		for i := range slice {
			slice[i] -= shift
			slice[i] *= scale
		}
	})
}

// Equalize equalizes the histogram
func (xs Xs) equalize() {
	xs.stretch()
	const a = 0xFFFF // discretization parameter
	buckets, cdf := make([]int, a), make(Xs, a)

	// count...
	for _, v := range xs {
		buckets[int(v*(a-1))]++
	}

	// ratio
	for i, v := range buckets {
		cdf[i] += cdf[clamp(i-1, len(cdf))]
		cdf[i] += X(v)
	}
	cdf.stretch()

	// transform
	for i, v := range xs {
		xs[i] = cdf[int(v*(a-1))]
	}

}

// Convlution kernel
type kernel struct {
	Xs
	n int
}

func (k kernel) String() string {
	return k.StringChop(k.n)
}

// GaussianKernel produces a 2-dimensional Gaussian kernel with sigma s
func GaussanKernel(n int) kernel {
	s := float64(n) / 6 // sigma
	k := kernel{make(Xs, n*n), n}
	for i := range k.Xs {
		x, y := i%n-n/2, i/n-n/2
		k.Xs[i] = X(math.Exp(-float64(x*x+y*y) / float64(2*s*s)))
	}

	k.stretch()
	return k
}

// LaplaceGaussKernel produces a 2-dimensional Gaussian kernel with sigma s; watchout, parameter is radius
func LaplaceGaussKernel(r int) kernel {
	n := r*2 + 1
	s := float64(n) / 6 // sigma
	k := kernel{make(Xs, n*n), n}

	for i := range k.Xs {
		x, y := i%n-n/2, i/n-n/2
		k.Xs[i] = -X(1 - float64(x*x+y*y)/float64(2*s*s))
		k.Xs[i] *= X(math.Exp(-float64(x*x+y*y) / float64(2*s*s)))
		k.Xs[i] /= X(math.Pi * math.Pow(s, 4))
	}

	return k
}

// Convolve transforms an image in-place(ish) by convolution with a given kernel.
func Convolve(img *image, k kernel) {
	out := make(Xs, len(img.Xs))
	img.Split(runtime.NumCPU(), func(slice Xs, offset int) {
		for i := range slice {
			j := offset + i
			x, y := xy(j, img.w)

			for h, v := range k.Xs {
				kx, ky := xyc(h, k.n)
				out[j] += v * img.at(x+kx, y+ky)
			}
		}
	})

	img.Xs = out
}
