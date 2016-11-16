package main

import (
	"os"
	"log"
	"path/filepath"
	"image"
	"image/color"
	"github.com/ilackarms/sprite-locator/models"
	"io/ioutil"
	"encoding/json"
	"image/png"
)

func main() {
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
			spritePixels := findSprite(x, y, img, bgColor, []image.Point{}, sprites)
			sprite := bounds(spritePixels)
			log.Printf("found a sprite with bounds %v", sprite)
			sprites = append(sprites, sprite)
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

func bounds(points []image.Point) image.Rectangle {
	var minX, minY, maxX, maxY int
	for _, point := range points {
		if minX < point.X {
			minX = point.X
		}
		if minY < point.Y {
			minY = point.Y
		}
		if maxX > point.X {
			maxX = point.X
		}
		if maxY > point.Y {
			maxY = point.Y
		}
	}
	return image.Rect(minX, minY, maxX, maxY)
}

func findSprite(x, y int, img image.Image, bgColor color.Color, points []image.Point, foundSprites []image.Rectangle) []image.Point {
	//out of bounds
	if x < img.Bounds().Min.X || x > img.Bounds().Max.X - 1 ||
		y < img.Bounds().Min.Y || y > img.Bounds().Max.Y - 1 {
		return points
	}
	//found a bg pixel
	if img.At(x, y) == bgColor {
		return points
	}
	//inspecting a pixel in a sprite already counted
	for _, sprite := range foundSprites {
		if image.Pt(x, y).In(sprite) {
			return points
		}
	}

	//found a point that we already marked
	for _, point := range points {
		if point.X == x && point.Y == y {
			return points
		}
	}
	//add x,y to points
	pixel := image.Point{X: x, Y: y}
	log.Printf("adding pixel %v", pixel)
	points = append(points, pixel)
	//recurse over N, E, S, W
	points = append(points, findSprite(x-1, y, img, bgColor, points, foundSprites)...)
	points = append(points, findSprite(x+1, y, img, bgColor, points, foundSprites)...)
	points = append(points, findSprite(x, y-1, img, bgColor, points, foundSprites)...)
	points = append(points, findSprite(x, y+1, img, bgColor, points, foundSprites)...)
	return points
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