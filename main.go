package main

import (
	"encoding/json"
	"github.com/ilackarms/sprite-locator/models"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"
	"strconv"
)

var markedPixels map[image.Point]bool
var margin int

func main() {
	margin = 4
	minImageHeight := 0
	if userMargin := os.Getenv("PIXEL_MARGIN"); userMargin != "" {
		usrM, err := strconv.Atoi(userMargin)
		if err != nil {
			log.Fatalf("%s is not a valid integer. unset PIXEL_MARGIN or give a valid value", userMargin)
		}
		margin = usrM
	}
	if userMinHeight := os.Getenv("MIN_IMAGE_HEIGHT"); userMinHeight != "" {
		usrM, err := strconv.Atoi(userMinHeight)
		if err != nil {
			log.Fatalf("%s is not a valid integer. unset MIN_IMAGE_HEIGHT or give a valid value", userMinHeight)
		}
		minImageHeight = usrM
	}

	markedPixels = make(map[image.Point]bool)
	args := os.Args
	if len(args) != 3 {
		log.Fatal("usage sprite-locator <filename> <out-file>")
	}
	inFile := args[1]
	outFile := args[2]
	path, err := filepath.Abs(inFile)
	if err != nil {
		log.Fatalf("abs path %v: %v", inFile, err)
	}
	log.Printf("reading image at %v", path)
	reader, err := os.Open(path)
	if err != nil {
		log.Fatalf("open %v: %v", path, err)
	}
	img, err := png.Decode(reader)
	if err != nil {
		log.Fatalf("reading err: %v", err)
	}

	log.Printf("finding sprites for dimensions %v", img.Bounds())

	bgColor := findBgColor(img)

	sprites := []image.Rectangle{}
	//mark all pixels that are not bgcolor
	scanImage(img, func(img image.Image, x, y int) {
		if img.At(x, y) != bgColor {
			sprite := newSpriteBounds()
			sprite.findBounds(x, y, img, bgColor, sprites)
			rect := image.Rectangle{sprite.min, sprite.max}
			if !rect.Empty() {
				if minImageHeight > 0 {
					if rect.Bounds().Size().Y < minImageHeight {
						log.Printf("sprite with bounds %v too small, ignoirng", rect)
						return
					}
				}
				sprites = append(sprites, rect)
				log.Printf("found a sprite with bounds %v; total sprites found: %v", rect, len(sprites))
				time.Sleep(time.Second * 10)
			}
		}
	})
	spriteSheet := models.Spritesheet{}

	for _, sprite := range sprites {
		spriteSheet.Sprites = append(spriteSheet.Sprites,
			models.Sprite{
				Min: models.Point{sprite.Min.X, sprite.Min.Y},
				Max: models.Point{sprite.Max.X, sprite.Max.Y},
			},
		)
	}
	data, err := json.Marshal(spriteSheet)
	if err != nil {
		log.Fatalf("marshalling sprite sheet metadata: %v", err)
	}

	if err := ioutil.WriteFile(outFile, data, 0644); err != nil {
		log.Fatalf("writing sprite sheet metadata: %v", err)
	}
	log.Printf("metadata sheet with %v sprites written to %s", len(spriteSheet.Sprites), outFile)
}

type spriteBounds struct {
	min, max image.Point
}

func newSpriteBounds() *spriteBounds {
	return &spriteBounds{
		min: image.Pt(math.MaxInt64, math.MaxInt64),
		max: image.Pt(-1, -1),
	}
}

func (cp *spriteBounds) findBounds(x, y int, img image.Image, bgColor color.Color, sprites []image.Rectangle) {
	pixel := image.Point{X: x, Y: y}
	//log.Printf("inspecting %v", pixel)
	//already inspected this pixel
	if marked := markedPixels[pixel]; marked {
		//log.Printf("REJECTED: %v,%v is used", x, y)
		return
	}
	markedPixels[pixel] = true

	//out of bounds
	if !pixel.In(img.Bounds()) {
		//log.Printf("REJECTED: %v,%v is out of bounds", x, y)
		return
	}
	//found a bg pixel
	if img.At(x, y) == bgColor {
		//log.Printf("REJECTED: %v,%v is bg", x, y)
		return
	}

	//inspecting a pixel in a sprite already counted
	for _, sprite := range sprites {
		if image.Pt(x, y).In(sprite) {
			//log.Printf("REJECTED: %v,%v is in a sprite already", x, y)
			return
		}
	}

	var out bool
	//resize bound
	if x <= cp.min.X {
		out = true
		cp.min.X = x
	}
	if y <= cp.min.Y {
		out = true
		cp.min.Y = y
	}
	if x >= cp.max.X {
		out = true
		cp.max.X = x
	}
	if y >= cp.max.Y {
		out = true
		cp.max.Y = y
	}

	if !out {
		//log.Printf("REJECTED: %v is inside bounds {%v:%v}", pixel, cp.min, cp.max)
		return
	}

	//log.Printf("adding point %v,%v", x, y)

	//recurse over left right up down pixels within a given margin
	for i := 1; i <= margin; i++ {
		cp.findBounds(x-i, y, img, bgColor, sprites)
		cp.findBounds(x+i, y, img, bgColor, sprites)
		cp.findBounds(x, y-i, img, bgColor, sprites)
		cp.findBounds(x, y+i, img, bgColor, sprites)
	}
	return
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
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X-1; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y-1; y++ {
			callback(img, x, y)
		}
	}
}
