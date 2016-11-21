package algorithm

import (
	"image"
	"image/color"
	"log"
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
	marked := make(map[image.Point]bool)
	scanImage(img, func(img image.Image, x, y int) {
		if img.At(x, y) != bgColor {
			sprite := newConnectedPixels(a.Margin, marked)
			sprite.findConnectingPixels(x, y, img, bgColor, sprites)
			rect := sprite.getBounds()
			if !rect.Empty() {
				if a.MinImageHeight > 0 {
					if rect.Bounds().Size().Y < a.MinImageHeight {
						//log.Printf("sprite with bounds %v too small, ignoirng", rect)
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

type connectedPixels struct {
	pixels []image.Point
	marked map[image.Point]bool
	margin int
}

func newConnectedPixels(margin int, marked map[image.Point]bool) *connectedPixels {
	return &connectedPixels{
		margin: margin,
		marked: marked,
	}
}

func (cp *connectedPixels) getBounds() image.Rectangle {
	if len(cp.pixels) < 2 {
		return image.Rect(0,0,0,0)
	}
	px0 := cp.pixels[0]
	minX := px0.X
	maxX := px0.X
	minY := px0.Y
	maxY := px0.Y
	for _, px := range cp.pixels {
		if px.X < minX {
			minX = px.X
		}
		if px.Y < minY {
			minY = px.Y
		}
		if px.X > maxX {
			maxX = px.X
		}
		if px.Y > maxY {
			maxY = px.Y
		}
	}
	return image.Rect(minX, minY, maxX, maxY)
}

func (cp *connectedPixels) findConnectingPixels(x, y int, img image.Image, bgColor color.Color, sprites []image.Rectangle) {
	pixel := image.Point{X: x, Y: y}
	//log.Printf("inspecting %v", pixel)
	//already inspected this pixel
	if marked := cp.marked[pixel]; marked {
		//log.Printf("REJECTED: %v,%v is used", x, y)
		return
	}
	cp.marked[pixel] = true

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

	cp.pixels = append(cp.pixels, pixel)

	//log.Printf("adding point %v,%v", x, y)

	//recurse over left right up down pixels within a given margin
	for i := 1; i <= cp.margin; i++ {
		cp.findConnectingPixels(x - i, y, img, bgColor, sprites)
		cp.findConnectingPixels(x + i, y, img, bgColor, sprites)
		cp.findConnectingPixels(x, y - i, img, bgColor, sprites)
		cp.findConnectingPixels(x, y + i, img, bgColor, sprites)
	}
	return
}
