package watermark

import (
	"andia/config"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"time"
)

var watermarkPath = ""
var logoPath = ""
var height = 0

func init() {
	rand.Seed(time.Now().UnixNano())
	config := config.Read()
	watermarkPath = config.Watermark.Watermark
	logoPath = config.Watermark.logo
	height = config.Watermark.Height
}

func getRandomInt(l int) int {
	return rand.Intn(l)
}

func getRandomColor() color.Color {
	return palette.WebSafe[getRandomInt(len(palette.WebSafe))]
}

func Watermark(imagePath string, output string) error {
	imgb, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	img, err := jpeg.Decode(imgb)
	if err != nil {
		return err
	}
	defer imgb.Close()

	bounds := img.Bounds()

	target := image.NewRGBA(bounds)
	draw.Draw(target, bounds, img, image.ZP, draw.Src)

	col := getRandomColor()

	for i := 0; i < bounds.Size().X; i++ {
		for j := 0; j < height; j++ {
			target.Set(i, bounds.Size().Y-j, col)
		}
	}

	wmb, err := os.Open(watermarkPath)
	if err != nil {
		return err
	}
	watermark, err := png.Decode(wmb)
	if err != nil {
		return err
	}
	defer wmb.Close()

	draw.Draw(target, bounds.Add(image.Point{bounds.Size().X/2 - watermark.Bounds().Size().X/2, bounds.Size().Y - watermark.Bounds().Size().Y}), watermark, image.ZP, draw.Over)

	lmp, err := os.Open(logoPath)
	if err != nil {
		return err
	}
	logo, err := png.Decode(lmp)
	if err != nil {
		return err
	}
	defer lmp.Close()

	for y := 10; y < bounds.Size().Y-100; y += logo.Bounds().Size().Y * 2 {
		draw.Draw(target, bounds.Add(image.Point{getRandomInt(bounds.Size().X), y}), logo, image.ZP, draw.Over)
	}

	out, err := os.Create(output)
	if err != nil {
		return err
	}

	err = jpeg.Encode(out, target, nil)
	if err != nil {
		return err
	}
	return nil
}
