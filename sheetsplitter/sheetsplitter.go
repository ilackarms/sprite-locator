package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
	"encoding/json"
	"image"
	"image/png"
	"os"
	"image/color"
	"github.com/golang/freetype"
	"log"
	"strings"
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
	imgFile := flag.String("img", "", "image file for drawing debugging boxes")
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
				x0 := subsheet.Start.X + col * (width + 1)
				box := Box{
					X: x0+1,
					Y: y0+1,
					W: width-1,
					H: height-2,
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
	if *imgFile != "" {
		must(drawDebugImage(*imgFile, atlas))
	}
}

func drawDebugImage(imgFile string, atlas Atlas) error {
	reader, err := os.Open(imgFile)
	if err != nil {
		return fmt.Errorf("open %v: %v", imgFile, err)
	}
	img, err := png.Decode(reader)
	if err != nil {
		return fmt.Errorf("reading err: %v", err)
	}
	newImage := image.NewRGBA(img.Bounds())
	white := color.RGBA{255,255,255,255}
	scanImage(img, func(img image.Image, x, y int) {
		newImage.Set(x, y, img.At(x, y))
		if _, _, _, a := img.At(x, y).RGBA(); a == 0 {
			newImage.Set(x, y, white)
		}
	})
	fontBytes, err := ioutil.ReadFile(os.Getenv("HOME")+"/workspace/scratch/fonts/Lato-Regular.ttf")
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
	colors := makeColors(atlas)
	for i, frame := range atlas.Frames {
		drawBox(newImage, frame.Box, colors, c, i)
	}

	outFile := strings.Replace(imgFile, ".png", ".debug.png", -1)

	//create or open file
	out, err := os.Create(outFile)
	if err != nil {
		log.Printf("WARN: creating file: %v", err)
		//open
		out, err = os.Open(outFile)
		if err != nil {
			return fmt.Errorf("opening existing file: %v", err)
		}
	}

	return png.Encode(out, newImage)
}

func drawBox(img *image.RGBA, box Box, colors []color.Color, context *freetype.Context, i int) {
	c := colors[i%len(colors)]
	for x := box.X; x < box.X + box.W; x++ {
		img.Set(x, box.Y, c)
		img.Set(x, box.Y+box.H, c)
	}
	for y := box.Y; y < box.Y + box.H; y++ {
		img.Set(box.X, y, c)
		img.Set(box.X+box.W, y, c)
	}
	context.DrawString(fmt.Sprintf("%v", i), freetype.Pt(box.X, box.Y))
}

func makeColors(atlas Atlas) []color.Color {
	boxColors := make([]color.Color, 256 * 256 * 256)
	i := 0
	for r := uint8(255); r >= 0; r-- {
		g := uint8(255 - r)
		b := uint8(0)
		boxColors[i] = color.RGBA{R: r, G: g, B: b, A: 255}
		i++
		if i > len(atlas.Frames) {
			break
		}
	}
	for g := uint8(255); g >= 0; g-- {
		b := uint8(255 - g)
		r := uint8(0)
		boxColors[i] = color.RGBA{R: r, G: g, B: b, A: 255}
		i++
		if i > len(atlas.Frames) {
			break
		}
	}
	for b := uint8(255); b >= 0; b-- {
		r := uint8(255 - b)
		g := uint8(0)
		boxColors[i] = color.RGBA{R: r, G: g, B: b, A: 255}
		i++
		if i > len(atlas.Frames) {
			break
		}

	}
	return boxColors
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