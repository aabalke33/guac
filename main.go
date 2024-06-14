package main

import (
	"guac/dev"
)

func main() {
	timer := timer.StartTimer()
	img, _ := openImage(url)
	timer.Lap("Opened Image")
	outputImage := getPixels(img)
	timer.Lap("Apply Filter")
	writeImage(outputImage)
	timer.Stop("Write Image")
}
