package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
	"encoding/json"
)

//generates an atlas directly from a single spritesheet
//diablo-formatted

func must(err interface{}) {
	if err != nil {
		logrus.Fatal(err)
	}
}

var rows = []string{
	"s",
	"sw",
	"w",
	"nw",
	"n",
	"ne",
	"e",
	"se",
}

func main() {
	metaFile := flag.String("meta", "", "metadata file that matches []subsheet format")
	flag.Parse()
	if *metaFile == "" {
		must("-meta must be set")
	}
	data, err := ioutil.ReadFile(*metaFile)
	must(err)
	var sheet Sheet
	err = yaml.Unmarshal(data, &sheet)
	must(err)

	var atlas Atlas
	//create atlas
	for _, subsheet := range sheet.Subsheets {
		animationName := subsheet.Name
		width := (subsheet.End.X - subsheet.Start.X)/ subsheet.Columns
		//edge case: subsheet only has a single row
		if subsheet.SingleRow {
			height := (subsheet.End.Y - subsheet.Start.Y)
			for col := 0; col < subsheet.Columns; col++ {
				//fill in every direction with the same column in the atlas
				for _, direction := range rows {
					frameName := fmt.Sprintf("%s.%s.%04d", animationName, direction, col)
					x0 := subsheet.Start.X + col * width
					y0 := subsheet.Start.Y
					box := Box{
						X: x0,
						Y: y0,
						W: width,
						H: height,
					}
					atlas.Frames = append(atlas.Frames, Frame{
						Filename: frameName,
						Box: box,
					})
				}
			}
			continue
		}
		height := (subsheet.End.Y - subsheet.Start.Y) / len(rows)
		for row, direction := range rows {
			y0 := subsheet.Start.Y + row * (height + sheet.RowSpacing)
			for col := 0; col < subsheet.Columns; col++ {
				frameName := fmt.Sprintf("%s.%s.%04d", animationName, direction, col)
				x0 := subsheet.Start.X + col * width
				box := Box{
					X: x0,
					Y: y0,
					W: width,
					H: height,
				}
				atlas.Frames = append(atlas.Frames, Frame{
					Filename: frameName,
					Box: box,
				})
			}
		}
	}
	data, err = json.Marshal(atlas)
	must(err)
	fmt.Printf("%s", data)
}
