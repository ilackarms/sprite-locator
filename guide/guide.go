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
	"sort"
	"image/color"
	"math"
	"github.com/golang/freetype"
)

var spriteMargin int

func main() {
	marginPtr := flag.Int("margin", 0, "margin around sprites for raytrace")
	outPtr := flag.String("out", "out", "output directory")
	jsonPtr := flag.String("json", "", "json file")
	imagePtr := flag.String("image", "", "image file")
	flag.Parse()
	spriteMargin = *marginPtr

	if *imagePtr == "" || *jsonPtr == "" {
		fmt.Println("usage: guide -image <image.png> -json <bounds.json> [-out <outdir>] [-margin int]")
		fmt.Printf("you gave me: %v\n", os.Args)
		os.Exit(-1)
	}
	if err := guide(*imagePtr, *jsonPtr, *outPtr); err != nil {
		log.Fatal(err)
	}
	log.Print("OK")
}

func guide(imgFile, jsonFile, outDir string) error {
	log.Printf("using: \n\timgFile: %v\n\tjsonFile: %v\n\toutDir: %v\n\tmargin %v", imgFile, jsonFile, outDir, spriteMargin)
	os.MkdirAll(outDir, 0755)

	path, err := filepath.Abs(imgFile)
	if err != nil {
		return errors.New(err, "abs path %v: %v", imgFile)
	}
	log.Printf("reading image at %v", path)
	reader, err := os.Open(path)
	if err != nil {
		return errors.New(err, "open %v", path)
	}
	img, err := png.Decode(reader)
	if err != nil {
		return errors.New(err, "reading err")
	}

	raw, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return errors.New(err, "reading json file")
	}
	var spritesheet models.Spritesheet
	if err := json.Unmarshal(raw, &spritesheet); err != nil {
		return errors.New(err, "failed to unmarshal spritesheet")
	}
	sortedSpritesheet := sortSheet(&spritesheet)
	if err := writeSheet(sortedSpritesheet, filepath.Join(outDir, jsonFile)); err != nil {
		return errors.New(err, "overwriting spritesheet")
	}
	return drawGuides(img, &sortedSpritesheet, filepath.Join(outDir, imgFile))
}

func sortSheet(sheet *models.Spritesheet) models.Spritesheet {
	sortedSprites := []models.Sprite{}
	min, max := getBounds(sheet)
	for len(sheet.Sprites) > 0 {
		log.Printf("finding top row of sprites in %v size sheet", len(sheet.Sprites))
		sortedSprites = append(sortedSprites, popTopRow(sheet, min, max)...)
		log.Printf("%v done, %v unsorted remaining", len(sortedSprites), len(sheet.Sprites))
	}
	return models.Spritesheet{Sprites: sortedSprites}
}

// like raycasting
// draw a line down from each x pixel
// when we hit a sprite, pop it from the sheet and add to the toprow
func popTopRow(sheet *models.Spritesheet, min, max image.Point) []models.Sprite {
	log.Printf("popping top row of bounds %v,%v", min, max)
	topRow := []models.Sprite{}
	notPopped := []models.Sprite{}
	//find topLeft Sprite
	var topLeft models.Sprite
	minDist := math.MaxFloat64
	for _, sprite := range sheet.Sprites {
		//dist = distance of min point from origin
		dist := math.Sqrt(math.Pow(float64(sprite.Min.X), 2)+math.Pow(float64(sprite.Min.Y), 2))
		if dist < minDist {
			minDist = dist
			topLeft = sprite
		}
	}
	topRow = append(topRow, topLeft)
	//draw a horizontal lince from center of topLeft
	c := center(topLeft)
	for _, sprite := range sheet.Sprites {
		//skip topleft
		if sprite == topLeft {
			continue
		}

		//draw a line from center
		for x := c.X; x <= max.X; x++{
			pt := image.Pt(x, c.Y)
			if inWithMargin(pt, sprite.Rect()) {
				topRow = append(topRow, sprite)
				break
			}
		}
	}
	for _, sprite := range sheet.Sprites {
		popped := false
		for _, poppedSprite := range topRow {
			if poppedSprite == sprite {
				popped = true
				break
			}
		}
		if !popped {
			notPopped = append(notPopped, sprite)
		}
	}
	log.Printf("%v popped, %v not popped", len(topRow), len(notPopped))
	//sort row by x position
	sorter := spriteSorter(topRow)
	sort.Sort(sorter)
	log.Printf("len sorter: %v, len topRow: %v", sorter.Len(), len(topRow))
	sheet.Sprites = notPopped
	return []models.Sprite(sorter)
}

func inWithMargin(p1 image.Point, rect image.Rectangle) bool {
	return p1.In(image.Rect(
		rect.Min.X-spriteMargin,
		rect.Min.Y-spriteMargin,
		rect.Max.X+spriteMargin,
		rect.Max.Y+spriteMargin,
	))
}

func dist(p1, p2 image.Point) float64 {
	return math.Sqrt(math.Pow(float64(p2.X-p1.X), 2)+math.Pow(float64(p2.Y-p1.Y), 2))
}

func getBounds(sheet *models.Spritesheet) (image.Point, image.Point) {
	minX := math.MaxInt64
	minY := sheet.Sprites[0].Min.Y
	maxX := sheet.Sprites[0].Max.X
	maxY := sheet.Sprites[0].Max.Y
	for _, sprite := range sheet.Sprites {
		if sprite.Min.X < minX {
			minX = sprite.Min.X
		}
		if sprite.Min.Y < minY {
			minY = sprite.Min.Y
		}

		if sprite.Max.X > maxX {
			maxX = sprite.Max.X
		}
		if sprite.Max.Y > maxY {
			maxY = sprite.Max.Y
		}
	}
	return image.Pt(minX, minY), image.Pt(maxX, maxY)
}

func center(sprite models.Sprite) image.Point {
	return image.Pt((sprite.Min.X+sprite.Max.X)/2, (sprite.Min.Y+sprite.Max.Y)/2)
}

type spriteSorter []models.Sprite
func (a spriteSorter) Len() int           { return len(a) }
func (a spriteSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a spriteSorter) Less(i, j int) bool { return center(a[i]).X < center(a[j]).X }

func writeSheet(sprites models.Spritesheet, outFile string) error {
	data, err := json.Marshal(sprites)
	if err != nil {
		return errors.New(err, "marshalling sprite sheet metadata")
	}

	if err := ioutil.WriteFile(outFile, data, 0644); err != nil {
		return errors.New(err, "writing sprite sheet metadata")
	}
	return nil
}

func drawGuides(img image.Image, sheet *models.Spritesheet, outFile string) error {
	log.Printf("drawing guides to %v", outFile)
	newImage := image.NewRGBA(img.Bounds())
	scanImage(img, func(img image.Image, x, y int) {
		newImage.Set(x, y, img.At(x, y))
	})
	boxColors := make([]color.Color, 256*256*256)
	i := 0
	for r := uint8(255); r >= 0; r-- {
		g := uint8(255 - r)
		b := uint8(0)
		boxColors[i] = color.RGBA{R: r, G: g, B: b, A: 255}
		i++
		if i > len(sheet.Sprites) {
			break
		}
	}
	for g := uint8(255); g >= 0; g-- {
		b := uint8(255 - g)
		r := uint8(0)
		boxColors[i] = color.RGBA{R: r, G: g, B: b, A: 255}
		i++
		if i > len(sheet.Sprites) {
			break
		}
	}
	for b := uint8(255); b >= 0; b-- {
		r := uint8(255 - b)
		g := uint8(0)
		boxColors[i] = color.RGBA{R: r, G: g, B: b, A: 255}
		i++
		if i > len(sheet.Sprites) {
			break
		}
	}

	fontBytes, err := ioutil.ReadFile("ANDRO.TTF")
	if err != nil {
		return err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	c := freetype.NewContext()
	c.SetDPI(72.0)
	c.SetFont(f)
	c.SetFontSize(24)
	c.SetClip(newImage.Bounds())
	c.SetDst(newImage)
	c.SetSrc(image.Black)

	for i, sprite := range sheet.Sprites {
		for _, pt := range boundingBoxPixels(sprite) {
			newImage.Set(pt.X, pt.Y, boxColors[i%len(sheet.Sprites)])
			drawNum(i, sprite.Min, c)
		}
	}

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

	return png.Encode(out, newImage)
}

func boundingBoxPixels(sprite models.Sprite) []image.Point {
	pixels := []image.Point{}
	for x := sprite.Min.X; x <= sprite.Max.X; x++ {
		//top line
		pixels = append(pixels, image.Pt(x, sprite.Min.Y))
		//bottom line
		pixels = append(pixels, image.Pt(x, sprite.Max.Y))
	}
	for y := sprite.Min.Y; y <= sprite.Max.Y; y++ {
		//left line
		pixels = append(pixels, image.Pt(sprite.Min.X, y))
		//right line
		pixels = append(pixels, image.Pt(sprite.Max.X, y))
	}
	return pixels
}

func drawNum(i int, loc models.Point, c *freetype.Context) {
	c.DrawString(fmt.Sprintf("%v", i), freetype.Pt(loc.X, loc.Y))
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