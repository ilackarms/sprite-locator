package main

import (
	"os"
	"fmt"
	"path/filepath"
	"log"
	"image/png"
	"io/ioutil"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/ilackarms/sprite-locator/models"
	"encoding/json"
	"image"
	"flag"
	"image/color"
	"math"
)

var spriteMargin int

func main() {
	//take in source image
	//take in boxes
	//scan boxes; find largest sprite
	//largest sprite becomes cell size
	//new image: cols = sqrt of len(boxes)
	//   width = cols * largest_sprite.width
	//   height = rows * largest_sprite.height

	//for each cell:
	//  locate center of sprite to draw
	// locate center of destination cell
	// translate = (c2 - c1)
	// apply translate to top left pixel, then draw

	imagePtr := flag.String("src", "", "image file")
	boxesPtr := flag.String("boxes", "", "boxes json file")
	outPtr := flag.String("src", "", "image file")
	flag.Parse()

	if *imagePtr == "" || *boxesPtr == "" || *outPtr == "" {
		fmt.Println("usage: sheetmaker -src <image.png> -boxes <boxes.json> -out <out.png>")
		fmt.Printf("you gave me: %v\n", os.Args)
		os.Exit(-1)
	}
	if err := makeSheet(*imagePtr, *boxesPtr, *outPtr); err != nil {
		log.Fatal(err)
	}
	log.Print("OK")
}

func makeSheet(imgFile, boxFile, outFile string) error {
	log.Printf("using: \n\timgFile: %v\n\boxFile: %v\n\toutDir: %v\n\tmargin %v", imgFile, boxFile, outFile, spriteMargin)

	path, err := filepath.Abs(imgFile)
	if err != nil {
		return errors.New(fmt.Sprintf("abs path %v", imgFile), err)
	}
	log.Printf("reading image at %v", path)
	reader, err := os.Open(path)
	if err != nil {
		return errors.New(fmt.Sprintf("open %v", path), err)
	}
	img, err := png.Decode(reader)
	if err != nil {
		return errors.New("reading err", err)
	}

	raw, err := ioutil.ReadFile(boxFile)
	if err != nil {
		return errors.New("reading box file", err)
	}
	var spriteSheet models.Spritesheet
	if err := json.Unmarshal(raw, &spriteSheet); err != nil {
		return errors.New(err, "failed to unmarshal spritesheet")
	}
	return drawNewSheet(img, &spriteSheet, outFile)
}

func drawNewSheet(img image.Image, sheet *models.Spritesheet, outFile string) error {
	//create or open file
	out, err := os.Create(outFile)
	if err != nil {
		log.Printf("WARN: creating file: %v", err)
		//open
		out, err = os.Open(outFile)
		if err != nil {
			return errors.New(err, "opening existing file")
		}
	}

	cellWidth, cellHeight := largestSpriteSize(sheet)
	cellCount := len(sheet.Sprites)
	cols := int(math.Sqrt(math.Floor(float64(cellCount))))
	rows := cellCount/cols + 1
	log.Printf("drawing new sheet to %v", outFile)
	log.Printf("cellWidth: %v, cellHeight: %v, rows: %v cols: %v", cellWidth, cellHeight, rows, cols)
	newImage := image.NewRGBA(image.Rect(0, 0, cellWidth * cols, cellHeight * rows))
	bgColor := color.RGBA{0,0,0,0}
	scanImage(img, func(img image.Image, x, y int) {
		newImage.Set(x, y, img.At(x, y))
		if _, _, _, a := img.At(x, y).RGBA(); a == 0{
			newImage.Set(x, y, bgColor)
		}
	})
	//draw each sprite from the original sprite sheet
	//into the corresponding cell on the new sheet
	for i, sprite := range sheet.Sprites {
		spriteCenter := image.Pt((sprite.Max.X+sprite.Min.X)/2, (sprite.Max.Y+sprite.Min.Y)/2)
		rowIndex := i/rows
		colIndex := i%cols
		cellStart := image.Pt(colIndex*cellWidth, rowIndex*cellHeight)
		cellCenter := image.Pt(cellStart+cellWidth/2, cellStart+cellHeight/2)
		offset := cellCenter.Sub(spriteCenter)
		for x := sprite.Min.X; x <= sprite.Max.X; x++ {
			for y := sprite.Min.Y; y <= sprite.Max.Y; y++ {
				px := img.At(x, y)
				newImage.Set(x + offset.X, y + offset.Y, px)
			}
		}
	}

	return png.Encode(out, newImage)
}

func largestSpriteSize(sheet *models.Spritesheet) (int, int) {
	var maxWidth, maxHeight int
	for _, sprite := range sheet.Sprites {
		if width:= (sprite.Max.X - sprite.Min.X); maxWidth < width {
			maxWidth = width
		}
		if height:= (sprite.Max.Y - sprite.Min.Y); maxHeight < height {
			maxHeight = height
		}
	}
	return maxWidth, maxHeight
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