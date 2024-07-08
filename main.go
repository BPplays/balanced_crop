package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"os"
	"time"

	g2bwebp "github.com/gen2brain/webp"
	"github.com/oliamb/cutter"
	"golang.org/x/image/webp"
)


func IsSimilar(c1, c2 color.NRGBA, SimilarityThreshold float64) bool {
	return math.Abs(float64(c1.R)-float64(c2.R)) <= SimilarityThreshold &&
		math.Abs(float64(c1.G)-float64(c2.G)) <= SimilarityThreshold &&
		math.Abs(float64(c1.B)-float64(c2.B)) <= SimilarityThreshold &&
		math.Abs(float64(c1.A)-float64(c2.A)) <= SimilarityThreshold
}


func crop_brd(img *image.Image, border_percent float64) *image.Image {
	// tlcol := img.At(0, 0)
	var SimilarityThreshold float64 = 54

	bounds := (*img).Bounds()
    width := bounds.Dx()
    height := bounds.Dy()

	short_exit_h := int(math.Max(float64(height) * 0.01, 5))

	long_exit_h := int(math.Max(float64(height) * 0.01, 5))

	if height < 20 {
		short_exit_h = 2
	}

	border_px := int(float64(height) * (border_percent / 100))

	var final_pixel_wcnt int = -1
	// var final_pixel_hcnt int = -1
	var wcnt_times int = 0
	var wcnt_times_long int = 0

	

	// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	// 	rightmostColor := (*img).At(bounds.Max.X-1, y).(color.NRGBA)
	// 	fmt.Printf("Pixel at (%d, %d) color: R=%d, G=%d, B=%d, A=%d\n", bounds.Max.X-1, y, rightmostColor.R, rightmostColor.G, rightmostColor.B, rightmostColor.A)
	// }

	
	for x := bounds.Min.X; x < width; x++ {
		tl_col := (*img).At(0, 0).(color.NRGBA)
		fmt.Println(IsSimilar(tl_col, tl_col, 10))


		wcnt_times_long = 0
		for y := bounds.Min.Y; y < height; y++ {
			if IsSimilar((*img).At(x, y).(color.NRGBA), tl_col, SimilarityThreshold) != true {
				final_pixel_wcnt = x
				wcnt_times++
				wcnt_times_long++
				fmt.Println((*img).At(x, y).(color.NRGBA))
			} else {
				wcnt_times = 0
			}
			if final_pixel_wcnt >= 0 && wcnt_times > short_exit_h && wcnt_times_long > long_exit_h {
				break
			}
		}


		wcnt_times_long = 0
		for y := height; y > bounds.Min.Y ; y-- {
			// fmt.Println(IsSimilar((*img).At(bounds.Max.X-1, y).(color.NRGBA), tl_col, SimilarityThreshold))
			// fmt.Println(final_pixel_wcnt, x)
			if IsSimilar((*img).At(width-x-1-1, y).(color.NRGBA), tl_col, SimilarityThreshold) != true {
				final_pixel_wcnt = x
				wcnt_times++
				wcnt_times_long++
				fmt.Println((*img).At(x, y).(color.NRGBA))
			} else {
				wcnt_times = 0
			}
			if final_pixel_wcnt >= 0 && wcnt_times > short_exit_h && wcnt_times_long > long_exit_h {
				break
			}
		}

		if final_pixel_wcnt >= 0 && wcnt_times > short_exit_h {
			fmt.Println(final_pixel_wcnt)
			break
		}

	}

	cwid := math.Min(float64(width - (final_pixel_wcnt - (border_px * 2)) * 2), float64(width))

	croppedImg, err := cutter.Crop(*img, cutter.Config{
		Width: int(math.Round(cwid)),
		Height: height,
		Mode: cutter.Centered,
	  })
	if err != nil {
		log.Fatalln("can't crop")
	}

	return &croppedImg

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


	file, err := os.Open(in)
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
		fmt.Println("NON-fatal error Dynamic lib file. encoding time will be slower:\n	", err)
		// return
	}

	cropstart := time.Now()
	croppedImg := crop_brd(&img, 0.5)
	fmt.Println("trim and crop time:", time.Since(cropstart))


	encstart := time.Now()
	// Create an output file
	outfile, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	g2bwebp.Encode(outfile, *croppedImg, g2bwebp.Options{Lossless: true, Method: 6, Exact: true})
	fmt.Println("encoding time:", time.Since(encstart))


}

func main() {
	read_crop("test.webp", "out.webp")
}