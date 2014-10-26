package main

// Threshold sets values below the input value to 0
func Threshold(img image, threshold X) {
	for i, v := range img.Xs {
		if v < threshold {
			img.Xs[i] = 0
		}
	}
}

// Count returns two vectors corresponding to the count of non-zero
// pixels in each row, column
func Count(img image) (nx, ny []float64) {
	nx, ny = make([]float64, img.w), make([]float64, img.h)
	for i, v := range img.Xs {
		if v != 0 {
			x, y := i%img.w, i/img.w
			nx[x]++
			ny[y]++
		}
	}
	return
}

// CoM: returns the Center of Mass (ie, the average non-zero point location)
// as two coordinates.
func CoM(img image) (cx, cy int) {
	tx, ty, tb := 0, 0, 0
	for i, v := range img.Xs {
		if v != 0 {
			x, y := i%img.w, i/img.w
			tb++
			tx += x
			ty += y
		}
	}

	cx, cy = tx/tb, ty/tb
	return
}

// Takes the count vectors (ie, from count) and produces an image
// corresponding to their product.
func Intersect(nx, ny []float64) image {
	sect := New(len(nx), len(ny))
	for i := range sect.Xs {
		x, y := i%sect.w, i/sect.w
		sect.Xs[i] = X(nx[x] * ny[y])
	}
	sect.stretch()
	return sect
}
