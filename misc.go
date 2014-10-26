package main

import (
	"errors"
	"fmt"
	eximg "image"
	"image/color"
	"image/png"
	"log"
	"os"
)

type lax interface{}

// Load reads a file at fname into a native image
func Load(fname string) image {
	f, err := os.Open(fname)
	must(err)
	defer f.Close()
	ximg, _, err := eximg.Decode(f)
	must(err)

	nimg, err := ToNative(ximg)
	must(err)

	return nimg
}

// ToNative takes an image.Image and gives us a version in the
// local.image format, normalized (ie, values in [0.0, 1.0]).
func ToNative(ximg eximg.Image) (image, error) {
	nimg := New(ximg.Bounds().Dx(), ximg.Bounds().Dy())
	switch ximg := ximg.(type) {
	case *eximg.Gray:
		for i := range ximg.Pix {
			nimg.Xs[i] = X(ximg.Pix[i])
		}
	case *eximg.RGBA:
		for i := 0; i < len(ximg.Pix); i += 4 {
			px := ximg.Pix[i : i+3] // TODO, currently ignoring alpha
			for j := range px {
				nimg.Xs[i/4] += X(ximg.Pix[i+j])
			}
		}
	default:
		return image{}, errors.New(fmt.Sprintf("Unsupported image type %T", ximg))
	}

	nimg.stretch()
	return nimg, nil
}

// Export takes a native image and writes it to file
func Export(nimg image, fname string) {
	f, err := os.Create(fname)
	must(err)
	defer f.Close()

	err = png.Encode(f, nimg)
	must(err)
}

func ck(e error) bool {
	if e != nil {
		log.Println(e)
		return true
	}
	return false
}

func must(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

// Clamp rounds input x up to 0 or down to top-1.
func clamp(x, top int) int {
	if x < 0 {
		return 0
	}
	if x >= top {
		return top - 1
	}

	return x
}

// maps i to xy coords
func xy(v, w int) (x, y int) {
	x, y = v%w, v/w
	return
}

// same as xy, but shifts to center
func xyc(v, w int) (x, y int) {
	x, y = v%w-w/2, v/w-w/2
	return
}

func (img image) ColorModel() color.Model {
	f := func(c color.Color) color.Color {
		r, g, b, _ := c.RGBA()
		return X(r+g+b) / 0xFFFF
	}
	return color.ModelFunc(f)
}
func (img image) Bounds() eximg.Rectangle {
	return eximg.Rect(0, 0, img.w, img.h)
}
func (img image) At(x, y int) color.Color {
	return img.Xs[x+img.w*y]
}
