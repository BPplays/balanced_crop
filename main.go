package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"time"

	g2bwebp "github.com/gen2brain/webp"
	"golang.org/x/image/webp"
)



func crop_brd(img *image.Image, border_percent float64) {
	// tlcol := img.At(0, 0)
	bounds := (*img).Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		rightmostColor := (*img).At(bounds.Max.X-1, y).(color.RGBA)
		fmt.Printf("Pixel at (%d, %d) color: R=%d, G=%d, B=%d, A=%d\n", bounds.Max.X-1, y, rightmostColor.R, rightmostColor.G, rightmostColor.B, rightmostColor.A)
	}
}


func imageToRGBA(src *image.Image) *image.RGBA {

    // No conversion needed if image is an *image.RGBA.
    if dst, ok := (*src).(*image.RGBA); ok {
        return dst
    }

    // Use the image/draw package to convert to *image.RGBA.
    b := (*src).Bounds()
    dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
    draw.Draw(dst, dst.Bounds(), (*src), b.Min, draw.Src)
    return dst
}

func read_crop(in string, out string) {
	var img image.Image
	var err error


	file, err := os.Open("test.webp")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()


	decstart := time.Now()
	// Decode the WebP file
	img, err = webp.Decode(file)
	if err != nil {
		fmt.Println("Error decoding WebP file:", err)
		return
	}
	fmt.Println("decoding time:", time.Since(decstart))

	err = g2bwebp.Dynamic()
	if err != nil {
		fmt.Println("NON-fatal error Dynamic lib file. encoding time will be slower:\n", err)
		// return
	}

	encstart := time.Now()
	// Create an output file
	outfile, err := os.Create("output.webp")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	g2bwebp.Encode(outfile, img, g2bwebp.Options{Lossless: true, Method: 6, Exact: true})
	fmt.Println("decoding time:", time.Since(encstart))


}

func main() {
	read_crop("test.webp", "out.webp")
}