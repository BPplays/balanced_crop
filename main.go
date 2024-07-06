package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	"golang.org/x/image/webp"
	"gopkg.in/gographics/imagick.v3/imagick"
)



func crop_brd(img *image.RGBA, border_percent float64) {
	// tlcol := img.At(0, 0)
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		rightmostColor := img.At(bounds.Max.X-1, y).(color.RGBA)
		fmt.Printf("Pixel at (%d, %d) color: R=%d, G=%d, B=%d, A=%d\n", bounds.Max.X-1, y, rightmostColor.R, rightmostColor.G, rightmostColor.B, rightmostColor.A)
	}
}

// rgbaToMagickWand converts an image.RGBA to a MagickWand
func rgbaToMagickWand(img *image.RGBA, mw *imagick.MagickWand) error {
    bounds := img.Bounds()
    width := bounds.Dx()
    height := bounds.Dy()

    // Create a new image in MagickWand
    if err := mw.NewImage(uint(width), uint(height), imagick.NewPixelWand()); err != nil {
        return err
    }

    // Set the image format
    if err := mw.SetImageFormat("RGBA"); err != nil {
        return err
    }

    // Import pixel data
    pixelData := img.Pix
    if err := mw.ImportImagePixels(0, 0, uint(width), uint(height), "RGBA", imagick.PIXEL_CHAR, pixelData); err != nil {
        return err
    }

    return nil
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
	var img image.RGBA

	file, err := os.Open("example.webp")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the WebP file
	imgw, err := webp.Decode(file)
	if err != nil {
		fmt.Println("Error decoding WebP file:", err)
		return
	}

	img = *(imageToRGBA(&imgw))

    // Create a new MagickWand
    mw := imagick.NewMagickWand()
    defer mw.Destroy()

    if err := rgbaToMagickWand(&img, mw); err != nil {
        log.Fatalf("Error converting RGBA to MagickWand: %v", err)
    }

    // Set WebP format
    mw.SetImageFormat("WEBP")

    // Enable lossless compression
    mw.SetOption("webp:lossless", "true")

    // Set the compression method to the highest effort (maximum value is 6)
    mw.SetOption("webp:method", "6")

    // Write the image
    if err := mw.WriteImage("output.webp"); err != nil {
        log.Fatalf("Error writing image: %v", err)
    }

    log.Println("Successfully converted image to lossless WebP with maximum compression effort")
}

func main() {
	imagick.Initialize()
    defer imagick.Terminate()
	read_crop("test.webp", "out.webp")
}