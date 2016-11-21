package algorithm

import "image"

type SpriteFindingAlgorithm interface {
	FindSprites(img image.Image) []image.Rectangle
}