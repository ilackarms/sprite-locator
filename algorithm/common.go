package algorithm

import (
	"image"
	"image/color"
)

type SpriteFindingAlgorithm interface {
	FindSprites(img image.Image) []image.Rectangle
}

func findBgColor(img image.Image) color.Color {
	//find most common color; this is background
	colorFrequencies := make(map[color.Color]int)
	scanImage(img, func(img image.Image, x, y int) {
		color := img.At(x, y)
		colorFrequencies[color] += 1
	})

	var bgColor color.Color
	maxFrequency := 0
	for color, frequency := range colorFrequencies {
		if frequency > maxFrequency {
			bgColor = color
			maxFrequency = frequency
		}
	}
	return bgColor
}

func scanImage(img image.Image, callback func(img image.Image, x, y int)) {
	// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
	// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y - 1; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X - 1; x++ {
			callback(img, x, y)
		}
	}
}
