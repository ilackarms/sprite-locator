package main

import (
	"encoding/json"
	"github.com/ilackarms/sprite-locator/algorithm"
	"github.com/ilackarms/sprite-locator/models"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"image"
	"fmt"
	"errors"
	"strings"
)

func main() {
	margin := 4
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
	var extractSprites bool
	if ex := os.Getenv("EXTRACT_SPRITES"); ex != "" && ex != "false" && ex != "0" {
		extractSprites = true
	}

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

	algorithm := algorithm.FloodFillAlgorithm{
		Margin:         margin,
		MinImageHeight: minImageHeight,
	}

	sprites := algorithm.FindSprites(img)

	spriteSheet := models.Spritesheet{}

	for i, sprite := range sprites {
		if extractSprites {
			fileName := fmt.Sprintf("%v_%v.png", strings.TrimSuffix(inFile, ".png"), i)
			if err := extractSprite(img, sprite, fileName); err != nil {
				log.Printf("ERROR: COULD NOT EXTRACT SPRITE: %v", err)
			}
		}

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

func extractSprite(srcImage image.Image, sprite image.Rectangle, outFile string) error {
	newImage := image.NewRGBA(srcImage.Bounds())

	// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
	// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
	for y := sprite.Min.Y; y < sprite.Max.Y - 1; y++ {
		for x := sprite.Min.X; x < sprite.Max.X - 1; x++ {
			newImage.Set(x, y, srcImage.At(x, y))
		}
	}

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