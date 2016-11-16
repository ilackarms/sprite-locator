package models

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

type Spritesheet struct {
	Sprites []Sprite `json:"sprites"`
}
