package main

import (
	"math"
	"runtime"
)

type X float64
type Xs []X

// implements color.Color
func (s X) RGBA() (r, g, b, a uint32) {
	x := uint32(s * 0xFFFF)
	return x, x, x, 0xFFFF
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

// GaussianKernel produces a 2-dimensional Gaussian kernel with sigma s
func GaussianKernel(s float64) kernel {
	n := int(math.Ceil(6 * s)) // sigma
	k := kernel{make(Xs, n*n), n}
	for i := range k.Xs {
		x, y := i%n-n/2, i/n-n/2
		k.Xs[i] = X(math.Exp(-float64(x*x+y*y) / float64(2*s*s)))
	}

	k.stretch()
	return k
}

// LaplaceGaussKernel produces a 2-dimensional Gaussian kernel with sigma s
func LaplaceGaussKernel(s float64) kernel {
	n := int(math.Ceil(6 * s)) // sigma
	k := kernel{make(Xs, n*n), n}

	for i := range k.Xs {
		x, y := i%n-n/2, i/n-n/2
		k.Xs[i] *= -X(float64(2*s*s) - float64(x*x+y*y))
		k.Xs[i] *= X(math.Exp(-float64(x*x+y*y) / float64(2*s*s)))
	}
	k.stretch()

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
				out[i] += v * img.at(x+kx, y+ky)
			}
		}
	})

	img.Xs = out
}
