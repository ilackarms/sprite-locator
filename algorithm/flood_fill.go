package algorithm

import (
	"image"
	"image/color"
	"log"
	"math"
)

type FloodFillAlgorithm struct {
	Margin         int
	MinImageHeight int
}

func (a *FloodFillAlgorithm) FindSprites(img image.Image) []image.Rectangle {
	bgColor := findBgColor(img)
	log.Printf("finding sprites in sheet %v with bg color %v", img.Bounds(), bgColor)
	sprites := []image.Rectangle{}
	//mark all pixels that are not bgcolor
	scanImage(img, func(img image.Image, x, y int) {
		if img.At(x, y) != bgColor {
			sprite := newSpriteBounds(a.Margin)
			sprite.findBounds(x, y, img, bgColor, sprites)
			rect := image.Rectangle{sprite.min, sprite.max}
			if !rect.Empty() {
				if a.MinImageHeight > 0 {
					if rect.Bounds().Size().Y < a.MinImageHeight {
						log.Printf("sprite with bounds %v too small, ignoirng", rect)
						return
					}
				}
				sprites = append(sprites, rect)
				log.Printf("found a sprite with bounds %v; total sprites found: %v", rect, len(sprites))
			}
		}
	})
	return sprites
}

type spriteBounds struct {
	min, max     image.Point
	markedPixels map[image.Point]bool
	margin       int
}

func newSpriteBounds(margin int) *spriteBounds {
	return &spriteBounds{
		min:          image.Pt(math.MaxInt64, math.MaxInt64),
		max:          image.Pt(-1, -1),
		margin:       margin,
		markedPixels: make(map[image.Point]bool),
	}
}

func (cp *spriteBounds) findBounds(x, y int, img image.Image, bgColor color.Color, sprites []image.Rectangle) {
	pixel := image.Point{X: x, Y: y}
	//log.Printf("inspecting %v", pixel)
	//already inspected this pixel
	if marked := cp.markedPixels[pixel]; marked {
		log.Printf("REJECTED: %v,%v is used", x, y)
		return
	}
	cp.markedPixels[pixel] = true

	//out of bounds
	if !pixel.In(img.Bounds()) {
		log.Printf("REJECTED: %v,%v is out of bounds", x, y)
		return
	}
	//found a bg pixel
	if img.At(x, y) == bgColor {
		log.Printf("REJECTED: %v,%v is bg", x, y)
		return
	}

	//inspecting a pixel in a sprite already counted
	for _, sprite := range sprites {
		if image.Pt(x, y).In(sprite) {
			log.Printf("REJECTED: %v,%v is in a sprite already", x, y)
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
		log.Printf("REJECTED: %v is inside bounds {%v:%v}", pixel, cp.min, cp.max)
		return
	}

	//log.Printf("adding point %v,%v", x, y)

	//recurse over left right up down pixels within a given margin
	for i := 1; i <= cp.margin; i++ {
		cp.findBounds(x - i, y, img, bgColor, sprites)
		cp.findBounds(x + i, y, img, bgColor, sprites)
		cp.findBounds(x, y - i, img, bgColor, sprites)
		cp.findBounds(x, y + i, img, bgColor, sprites)
	}
	return
}
