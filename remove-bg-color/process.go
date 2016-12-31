package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		log.Print("usage: process in.png out.png")
		log.Fatalf("args given: %v", args)
	}
	inFile := args[1]
	if err := process(inFile, inFile); err != nil {
		log.Fatal(err)
	}
	log.Print("OK")
}

func process(inFile, outFile string) error {
	path, err := filepath.Abs(inFile)
	if err != nil {
		return errors.New(fmt.Sprintf("abs path %v: %v", inFile, err))
	}
	log.Printf("reading image at %v", path)
	reader, err := os.Open(path)
	if err != nil {
		return errors.New(fmt.Sprintf("open %v: %v", path, err))
	}
	srcImage, err := png.Decode(reader)
	if err != nil {
		return errors.New(fmt.Sprintf("reading err: %v", err))
	}
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	bgColor := color.RGBA{R: 255, G: 255, B: 255, A: 0}
	newImage := image.NewRGBA(srcImage.Bounds())
	scanImage(srcImage, func(img image.Image, x, y int) {
		pixel := img.At(x, y)
		if equal(pixel, white) {
			pixel = bgColor
		}
		newImage.Set(x, y, pixel)
	})

	//create or open file
	out, err := os.Create(outFile)
	if err != nil {
		log.Printf("WARN: creating file: %v", err)
		//open
		out, err = os.Open(outFile)
		if err != nil {
			return errors.New(fmt.Sprintf("opening existing file: %v", err))
		}
	}

	return png.Encode(out, newImage)
}

func equal(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

func scanImage(img image.Image, callback func(img image.Image, x, y int)) {
	// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
	// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X-1; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y-1; y++ {
			callback(img, x, y)
		}
	}
}
