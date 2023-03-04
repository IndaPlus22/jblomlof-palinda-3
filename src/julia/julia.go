// Stefan Nilsson 2013-02-27

/*
FASTEST TIME PRE OPTIMIZATION TOOK: 17.9 seconds on my hardware (some shitty second hand windows pc) (using windows: Measure-Command { go run julia.go}).
LETS DO SOME MATHS, if the optimal time per routince is 1 ms. we want to divide into 18 goroutines, approx
					so easizest would just to be do let each picture get a go routine which then divides into two goroutines
					ofc this assumes each pixel takes the same amount of time

FASTEST TIME POST OPTIMIZATION TOOK: 8.7 seconds (using: Measure-Command { go run julia.go})
*/

// This program creates pictures of Julia sets (en.wikipedia.org/wiki/Julia_set).
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"os"
	"strconv"
	"sync"
)

type ComplexFunc func(complex128) complex128

var Funcs []ComplexFunc = []ComplexFunc{
	func(z complex128) complex128 { return z*z - 0.61803398875 },
	func(z complex128) complex128 { return z*z + complex(0, 1) },
	func(z complex128) complex128 { return z*z + complex(-0.835, -0.2321) },
	func(z complex128) complex128 { return z*z + complex(0.45, 0.1428) },
	func(z complex128) complex128 { return z*z*z + 0.400 },
	func(z complex128) complex128 { return cmplx.Exp(z*z*z) - 0.621 },
	func(z complex128) complex128 { return (z*z+z)/cmplx.Log(z) + complex(0.268, 0.060) },
	func(z complex128) complex128 { return cmplx.Sqrt(cmplx.Sinh(z*z)) + complex(0.065, 0.122) },
}

func main() {
	wgOuter := new(sync.WaitGroup)
	wgOuter.Add(len(Funcs))
	for n, fn := range Funcs {
		go CreatePng("picture-"+strconv.Itoa(n)+".png", fn, 1024, wgOuter)

	}
	wgOuter.Wait()
}

// CreatePng creates a PNG picture file with a Julia image of size n x n.
func CreatePng(filename string, f ComplexFunc, n int, wgOuter *sync.WaitGroup) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	err = png.Encode(file, Julia(f, n))
	if err != nil {
		log.Fatal(err)
	}

	wgOuter.Done()
}

// Julia returns an image of size n x n of the Julia set for f.
func Julia(f ComplexFunc, n int) image.Image {
	bounds := image.Rect(-n/2, -n/2, n/2, n/2)
	img := image.NewRGBA(bounds)
	wg := new(sync.WaitGroup)
	s := float64(n / 4)
	minX := bounds.Min.X
	maxX := bounds.Max.X
	// calling 3 goroutines wasn't faster
	// if minX this is completly useless. But we getting the boundires by ((maxX - minX) + minX)/ 3 or similiar
	//firstThird = (2*minX + maxX) / 3
	//secondThird = (minX + 2*maxX) / 3
	halfX := (minX + maxX) / 2
	wg.Add(2)

	go setPixels(f, img, minX, halfX, bounds.Min.Y, bounds.Max.Y, s, wg)
	go setPixels(f, img, halfX, maxX, bounds.Min.Y, bounds.Max.Y, s, wg)
	/*for i := bounds.Min.X; i < bounds.Max.X; i++ {
		for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
			n := Iterate(f, complex(float64(i)/s, float64(j)/s), 256)
			r := uint8(0)
			g := uint8(0)
			b := uint8(n % 32 * 8)
			img.Set(i, j, color.RGBA{r, g, b, 255})
		}
	}*/
	wg.Wait()
	return img
}

func setPixels(f ComplexFunc, img *image.RGBA, lowX, highX, lowY, highY int, s float64, wg *sync.WaitGroup) {
	for i := lowX; i < highX; i++ {
		for j := lowY; j < highY; j++ {
			n := Iterate(f, complex(float64(i)/s, float64(j)/s), 256)
			r := uint8(0)
			g := uint8(0)
			b := uint8(n % 32 * 8)
			img.Set(i, j, color.RGBA{r, g, b, 255})
		}
	}
	wg.Done()
}

// Iterate sets z_0 = z, and repeatedly computes z_n = f(z_{n-1}), n â‰¥ 1,
// until |z_n| > 2  or n = max and returns this n.
func Iterate(f ComplexFunc, z complex128, max int) (n int) {
	for ; n < max; n++ {
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			break
		}
		z = f(z)
	}
	return
}
