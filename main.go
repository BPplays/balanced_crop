package main

import (
	"bufio"
	"image"
	"log"

	"github.com/tidbyt/go-libwebp/test/util"
	"github.com/tidbyt/go-libwebp/webp"
)



func crop_brd(img *image.RGBA, border_percent float64) {
	tlcol := img.At(0, 0)
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
	var err error

	// Read binary data
	data := util.ReadFile("cosmos.webp")

	// Decode
	options := &webp.DecoderOptions{}
	img, err := webp.DecodeRGBA(data, options)
	if err != nil {
		panic(err)
	}

	img := util.ReadPNG("cosmos.png")

	// Create file and buffered writer
	io := util.CreateFile("encoded_cosmos.webp")
	w := bufio.NewWriter(io)
	defer func() {
		w.Flush()
		io.Close()
	}()

	config := webp.ConfigPreset(webp.PresetDefault, 90)

	// Encode into WebP
	if err := webp.EncodeRGBA(w, img.(*image.RGBA), config); err != nil {
		panic(err)
	}
}