package models

import "image"

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Sprite struct {
	//Upper left pixel
	Min Point `json:"min"`
	//Lower Right pixel
	Max Point `json:"max"`
}

func (s Sprite) Rect() image.Rectangle {
	return image.Rect(s.Min.X, s.Min.Y, s.Max.X, s.Max.Y)
}

type Spritesheet struct {
	Sprites []Sprite `json:"sprites"`
}
