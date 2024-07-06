package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/tidbyt/go-libwebp/test/util"
	"github.com/tidbyt/go-libwebp/webp"
)



func crop_brd(img *image.RGBA, border_percent float64) {
	// tlcol := img.At(0, 0)
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		rightmostColor := img.At(bounds.Max.X-1, y).(color.RGBA)
		fmt.Printf("Pixel at (%d, %d) color: R=%d, G=%d, B=%d, A=%d\n", bounds.Max.X-1, y, rightmostColor.R, rightmostColor.G, rightmostColor.B, rightmostColor.A)
	}
}


func read_crop(in string, out string) {
	var err error

	// Read binary data
	data := util.ReadFile(in)

	// Decode
	options := &webp.DecoderOptions{UseThreads: true}
	img, err := webp.DecodeRGBA(data, options)
	if err != nil {
		panic(err)
	}

	// Create file and buffered writer
	io := util.CreateFile(out)
	w := bufio.NewWriter(io)
	defer func() {
		w.Flush()
		io.Close()
	}()

	config, err := webp.ConfigLosslessPreset(9)
	if err != nil {
		log.Fatalln(err)
	}

	// Encode into WebP
	if err := webp.EncodeRGBA(w, img, config); err != nil {
		panic(err)
	}
}

func main() {
	read_crop("test.webp", "out.webp")
}