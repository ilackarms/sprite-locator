package main

//custom stuff
type Point struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

type Subsheet struct {
	Name       string `yaml:"name"`
	Start      Point `yaml:"start"`
	End        Point `yaml:"end"`
	Columns    int `yaml:"columns"`
	SingleRow  bool `yaml:"single_row,omitempty"`
	Reversed  bool `yaml:"reversed,omitempty"`
}

type Sheet struct {
	Subsheets []Subsheet `yaml:"subsheets"`
	RowSpacing int `yaml:"row_spacing"`
}


//Atlas stuff
type Atlas struct {
	Frames []Frame `json:"frames"`
}

type Frame struct {
	Filename string `json:"filename"`
	Box Box `json:"frame"`
	Rotated bool `json:"rotated"`
	Trimmed bool `json:"trimmed"`
	SpriteSourceSize struct {
			 X int `json:"x"`
			 Y int `json:"y"`
			 W int `json:"w"`
			 H int `json:"h"`
		 } `json:"spriteSourceSize"`
	SourceSize struct {
			 W int `json:"w"`
			 H int `json:"h"`
		 } `json:"sourceSize"`
}

type Box struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type Anims struct {
	Attack struct {
		       S []string `json:"s"`
		       Sw []string `json:"sw"`
		       W []string `json:"w"`
		       Nw []string `json:"nw"`
		       N []string `json:"n"`
		       Ne []string `json:"ne"`
		       E []string `json:"e"`
		       Se []string `json:"se"`
	       } `json:"attack"`
	Idle struct {
		       S []string `json:"s"`
		       Sw []string `json:"sw"`
		       W []string `json:"w"`
		       Nw []string `json:"nw"`
		       N []string `json:"n"`
		       Ne []string `json:"ne"`
		       E []string `json:"e"`
		       Se []string `json:"se"`
	       } `json:"idle"`
	Walk struct {
		       S []string `json:"s"`
		       Sw []string `json:"sw"`
		       W []string `json:"w"`
		       Nw []string `json:"nw"`
		       N []string `json:"n"`
		       Ne []string `json:"ne"`
		       E []string `json:"e"`
		       Se []string `json:"se"`
	       } `json:"walk"`
	GetHit struct {
		       S []string `json:"s"`
		       Sw []string `json:"sw"`
		       W []string `json:"w"`
		       Nw []string `json:"nw"`
		       N []string `json:"n"`
		       Ne []string `json:"ne"`
		       E []string `json:"e"`
		       Se []string `json:"se"`
	       } `json:"get_hit"`
	Die struct {
		       S []string `json:"s"`
		       Sw []string `json:"sw"`
		       W []string `json:"w"`
		       Nw []string `json:"nw"`
		       N []string `json:"n"`
		       Ne []string `json:"ne"`
		       E []string `json:"e"`
		       Se []string `json:"se"`
	       } `json:"die"`
	Spell struct {
		       S []string `json:"s"`
		       Sw []string `json:"sw"`
		       W []string `json:"w"`
		       Nw []string `json:"nw"`
		       N []string `json:"n"`
		       Ne []string `json:"ne"`
		       E []string `json:"e"`
		       Se []string `json:"se"`
	       } `json:"spell"`

}
