package main

import (
	"os"
	"log"
	"io/ioutil"
	"github.com/ilackarms/sprite-locator/models"
	"encoding/json"
	"strings"
	"strconv"
	"fmt"
)

func main(){
	boxFile := os.Args[1]
	animFile := os.Args[2]
	boxData, err := ioutil.ReadFile(boxFile)
	must(err)
	animData, err := ioutil.ReadFile(animFile)
	must(err)
	var boxes models.Spritesheet
	must(json.Unmarshal(boxData, &boxes))
	var anims Anims
	must(json.Unmarshal(animData, &anims))
	var atlas Atlas
	addRange(&atlas, boxes, "Attack.S", anims.Attack.S)
	addRange(&atlas, boxes, "Attack.Sw", anims.Attack.Sw)
	addRange(&atlas, boxes, "Attack.W", anims.Attack.W)
	addRange(&atlas, boxes, "Attack.Nw", anims.Attack.Nw)
	addRange(&atlas, boxes, "Attack.N", anims.Attack.N)
	addRange(&atlas, boxes, "Attack.Ne", anims.Attack.Ne)
	addRange(&atlas, boxes, "Attack.E", anims.Attack.E)
	addRange(&atlas, boxes, "Attack.Se", anims.Attack.Se)

	addRange(&atlas, boxes, "Die.S", anims.Die.S)
	addRange(&atlas, boxes, "Die.Sw", anims.Die.Sw)
	addRange(&atlas, boxes, "Die.W", anims.Die.W)
	addRange(&atlas, boxes, "Die.Nw", anims.Die.Nw)
	addRange(&atlas, boxes, "Die.N", anims.Die.N)
	addRange(&atlas, boxes, "Die.Ne", anims.Die.Ne)
	addRange(&atlas, boxes, "Die.E", anims.Die.E)
	addRange(&atlas, boxes, "Die.Se", anims.Die.Se)

	addRange(&atlas, boxes, "GetHit.S", anims.GetHit.S)
	addRange(&atlas, boxes, "GetHit.Sw", anims.GetHit.Sw)
	addRange(&atlas, boxes, "GetHit.W", anims.GetHit.W)
	addRange(&atlas, boxes, "GetHit.Nw", anims.GetHit.Nw)
	addRange(&atlas, boxes, "GetHit.N", anims.GetHit.N)
	addRange(&atlas, boxes, "GetHit.Ne", anims.GetHit.Ne)
	addRange(&atlas, boxes, "GetHit.E", anims.GetHit.E)
	addRange(&atlas, boxes, "GetHit.Se", anims.GetHit.Se)

	addRange(&atlas, boxes, "Idle.S", anims.Idle.S)
	addRange(&atlas, boxes, "Idle.Sw", anims.Idle.Sw)
	addRange(&atlas, boxes, "Idle.W", anims.Idle.W)
	addRange(&atlas, boxes, "Idle.Nw", anims.Idle.Nw)
	addRange(&atlas, boxes, "Idle.N", anims.Idle.N)
	addRange(&atlas, boxes, "Idle.Ne", anims.Idle.Ne)
	addRange(&atlas, boxes, "Idle.E", anims.Idle.E)
	addRange(&atlas, boxes, "Idle.Se", anims.Idle.Se)

	addRange(&atlas, boxes, "Spell.S", anims.Spell.S)
	addRange(&atlas, boxes, "Spell.Sw", anims.Spell.Sw)
	addRange(&atlas, boxes, "Spell.W", anims.Spell.W)
	addRange(&atlas, boxes, "Spell.Nw", anims.Spell.Nw)
	addRange(&atlas, boxes, "Spell.N", anims.Spell.N)
	addRange(&atlas, boxes, "Spell.Ne", anims.Spell.Ne)
	addRange(&atlas, boxes, "Spell.E", anims.Spell.E)
	addRange(&atlas, boxes, "Spell.Se", anims.Spell.Se)

	addRange(&atlas, boxes, "Walk.S", anims.Walk.S)
	addRange(&atlas, boxes, "Walk.Sw", anims.Walk.Sw)
	addRange(&atlas, boxes, "Walk.W", anims.Walk.W)
	addRange(&atlas, boxes, "Walk.Nw", anims.Walk.Nw)
	addRange(&atlas, boxes, "Walk.N", anims.Walk.N)
	addRange(&atlas, boxes, "Walk.Ne", anims.Walk.Ne)
	addRange(&atlas, boxes, "Walk.E", anims.Walk.E)
	addRange(&atlas, boxes, "Walk.Se", anims.Walk.Se)

	log.Printf("%+v", atlas)

	raw, err := json.Marshal(atlas)
	must(err)
	fmt.Print(string(raw))
}

func addRange(atlas *Atlas, boxes models.Spritesheet, animationName string, frameRange []string) {
	frameCount := 1
	for _, r := range frameRange {
		if strings.Contains(r, "..") {
			split := strings.Split(r, "..")
			b := split[0]
			e := split[1]
			begin, err := strconv.Atoi(b)
			must(err)
			end, err := strconv.Atoi(e)
			must(err)
			for i := begin; i <= end; i++ {
				frameName := fmt.Sprintf("%s%04d", animationName, frameCount)
				frame := getFrame(frameName, boxes, i)
				atlas.Frames = append(atlas.Frames, frame)
				frameCount++
			}
		} else {
			i, err := strconv.Atoi(r)
			must(err)
			frameName := fmt.Sprintf("%s%04d", animationName, frameCount)
			frame := getFrame(frameName, boxes, i)
			atlas.Frames = append(atlas.Frames, frame)
			frameCount++
		}
	}
}

func getFrame(frameName string, boxes models.Spritesheet, i int) Frame {
	box := boxes.Sprites[i]
	x0 := box.Min.X
	y0 := box.Min.Y
	w := box.Max.X - x0
	h := box.Max.Y - y0

	return Frame{
		Filename: frameName,
		Box: Box{
			X: x0,
			Y: y0,
			W: w,
			H: h,
		},
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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