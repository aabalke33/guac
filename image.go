package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"sync"
)

const MAX_BIT = float64(0xFFFF)

func openImage(url string) (img image.Image, err error) {
	f, err := os.Open(url)
	if err != nil {
		panic("Could not open file")
	}

	img, format, err := image.Decode(f)
	if err != nil {
		panic("Could not decode file")
	}

	if format != "jpeg" {
		panic("Image is not Jpeg")
	}

	return img, nil
}

func getPixels(img image.Image) *image.RGBA {

	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	outputImage := image.NewRGBA(rect)

	var wg sync.WaitGroup
	wg.Add(size.X)

	for x := 0; x < size.X; x++ {
		go func(i int) {
			defer wg.Done()
			for j := 0; j < size.Y; j++ {
				outputImage.Set(i, j, applyFilterChannel(i, j, img.At(i, j)))
			}
		}(x)

	}
	wg.Wait()
	return outputImage
}

func rgbToHsl(rInt, gInt, bInt uint16) (hue int, sat, lum float64) {

	r := float64(rInt) / MAX_BIT
	g := float64(gInt) / MAX_BIT
	b := float64(bInt) / MAX_BIT

	minimum := min(r, g, b)
	maximum := max(r, g, b)
	delta := (maximum - minimum)
	lum = (maximum + minimum) / 2

	if delta == 0 {
		return 0, 0, lum
	}

	sat = delta / (1 - math.Abs((2*lum)-1))
	if sat > 1 {
		sat = 1
	}

	var hueFloat float64
	switch maximum {
	case r:
		hueFloat = (g - b) / delta
	case g:
		hueFloat = (b-r)/delta + 2
	case b:
		hueFloat = (r-g)/delta + 4
	}

	hue = int(hueFloat * 60)

	if hue < 0 {
		hue = 360
	}

	return hue, sat, lum
}

func hslToRgb(hueInt int, sat, lum float64) (r, g, b uint16) {
	if sat > 1 {
		sat = 1
	}
	if sat < 0 {
		sat = 0
	}
	if hueInt > 360 {
		hueInt = hueInt - 360
	}
	if hueInt < 0 {
		hueInt = hueInt + 360
	}
	if lum > 1 {
		lum = 1
	}
	if lum < 0 {
		lum = 0
		sat = 0
	}
	//	sat = max(sat, 0)
	//	sat = min(sat, 1)
	//	lum = max(lum, 0)
	//	lum = min(lum, 1)
	//
	//	if hueInt > 360 {
	//		hueInt = hueInt - 360
	//	}
	//	if hueInt < 0 {
	//		hueInt = hueInt + 360
	//	}
	hue := float64(hueInt) / 360

	if sat == 0 {
		r = uint16(lum * MAX_BIT)
		g = uint16(lum * MAX_BIT)
		b = uint16(lum * MAX_BIT)
		return
	}

	c := (1 - math.Abs((2*lum)-1)) * sat
	x := c * (1 - math.Abs(math.Mod(hue*6, 2)-1))
	m := lum - (c / 2)

	var rP, gP, bP float64
	switch {
	case hueInt < 60:
		rP, gP, bP = c, x, 0
	case hueInt < 120:
		rP, gP, bP = x, c, 0
	case hueInt < 180:
		rP, gP, bP = 0, c, x
	case hueInt < 240:
		rP, gP, bP = 0, x, c
	case hueInt < 300:
		rP, gP, bP = x, 0, c
	case hueInt <= 360:
		rP, gP, bP = c, 0, x
	}

	r = uint16((rP + m) * MAX_BIT)
	g = uint16((gP + m) * MAX_BIT)
	b = uint16((bP + m) * MAX_BIT)
	return r, g, b
}

func applyFilter(pixel color.Color) (outputPixel color.Color) {

	r, g, b, a := pixel.RGBA()
	R := uint16(r)
	G := uint16(g)
	B := uint16(b)
	A := uint16(a)

	hue, sat, lum := rgbToHsl(R, G, B)
	//hue = hue + 10
	//sat = sat-0.5
	oR, oG, oB := hslToRgb(hue, sat, lum)

	return color.NRGBA64{oR, oG, oB, A}
}

func applyFilterChannel(i int, j int, pixel color.Color) (outputPixel color.Color) {

	r, g, b, a := pixel.RGBA()
	R := uint16(r)
	G := uint16(g)
	B := uint16(b)
	A := uint16(a)

	Hue, Sat, Lum := rgbToHsl(R, G, B)
	rHue, rSat, rLum := rgbToHsl(R, G, B)
	gHue, gSat, gLum := rgbToHsl(G, B, R)
	bHue, bSat, bLum := rgbToHsl(B, R, G)
    fmt.Printf("%d, %f ,%f\n", Hue, Sat, Lum)
    fmt.Printf("%d, %f ,%f\n", rHue, rSat, rLum)
    fmt.Printf("%d, %f ,%f\n", gHue, gSat, gLum)
    fmt.Printf("%d, %f ,%f\n", bHue, bSat, bLum)
	//hue = hue + 10
	//sat = sat-0.5
	oR, oG, oB := hslToRgb(Hue, Sat, Lum)

	return color.NRGBA64{oR, oG, oB, A}
}

func writeImage(outputImage *image.RGBA) {

	f, err := os.Create("./output.jpg")
	if err != nil {
		fmt.Println("Creating file:", err)
	}

	defer f.Close()
	if err = jpeg.Encode(f, outputImage, nil); err != nil {
		panic(err)
	}

}
