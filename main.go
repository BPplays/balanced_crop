package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	g2bwebp "github.com/gen2brain/webp"
	"github.com/oliamb/cutter"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/image/webp"
)


func IsSimilar(c1, c2 color.NRGBA, SimilarityThreshold float64) bool {
	return math.Abs(float64(c1.R)-float64(c2.R)) <= SimilarityThreshold &&
		math.Abs(float64(c1.G)-float64(c2.G)) <= SimilarityThreshold &&
		math.Abs(float64(c1.B)-float64(c2.B)) <= SimilarityThreshold &&
		math.Abs(float64(c1.A)-float64(c2.A)) <= SimilarityThreshold
}


func crop_brd_w(img *image.Image, border_percent *float64, SimilarityThreshold *float64, short_exit_mul *float64, long_exit_mul *float64) (*float64, *int) {
	bounds := (*img).Bounds()
    width := bounds.Dx()
    height := bounds.Dy()

	short_exit := int(math.Max(float64(width) * (*short_exit_mul), 5))

	long_exit := int(math.Max(float64(width) * (*long_exit_mul), 5))

	if width < 20 {
		short_exit = 2
	}

	border_px_wid := int(float64(width) * (*border_percent / 100))

	var final_pixel_cnt int = -1

	var wcnt_times int = 0
	var wcnt_times_long int = 0

	

	// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	// 	rightmostColor := (*img).At(bounds.Max.X-1, y).(color.NRGBA)
	// 	fmt.Printf("Pixel at (%d, %d) color: R=%d, G=%d, B=%d, A=%d\n", bounds.Max.X-1, y, rightmostColor.R, rightmostColor.G, rightmostColor.B, rightmostColor.A)
	// }

	
	for x := bounds.Min.X; x < width; x++ {
		tl_col := (*img).At(bounds.Min.X, bounds.Min.Y).(color.NRGBA)
		// fmt.Println(IsSimilar(tl_col, tl_col, 10))


		wcnt_times_long = 0
		for y := bounds.Min.Y; y < height; y++ {
			if IsSimilar((*img).At(x, y).(color.NRGBA), tl_col, *SimilarityThreshold) != true {
				final_pixel_cnt = x
				wcnt_times++
				wcnt_times_long++
				// fmt.Println((*img).At(x, y).(color.NRGBA))
			} else {
				wcnt_times = 0
			}
			if final_pixel_cnt >= 0 && (wcnt_times > short_exit || wcnt_times_long > long_exit) {
				break
			}
		}


		wcnt_times_long = 0
		for y := width; y > bounds.Min.Y ; y-- {
			// fmt.Println(IsSimilar((*img).At(bounds.Max.X-1, y).(color.NRGBA), tl_col, SimilarityThreshold))
			// fmt.Println(final_pixel_wcnt, x)
			if IsSimilar((*img).At(width-x-1, y).(color.NRGBA), tl_col, *SimilarityThreshold) != true {
				final_pixel_cnt = x
				wcnt_times++
				wcnt_times_long++
				// fmt.Println((*img).At(x, y).(color.NRGBA))
			} else {
				wcnt_times = 0
			}
			if final_pixel_cnt >= 0 && (wcnt_times > short_exit || wcnt_times_long > long_exit) {
				break
			}
		}

		// final_pixel_cnt = int(math.Min(float64(final_pixel_cnt1), float64(final_pixel_cnt2)))
		// if final_pixel_cnt < 0 {
		// 	if final_pixel_cnt1 >= 0 || final_pixel_cnt2 >= 0 {
		// 		final_pixel_cnt = 0
		// 	}
		// }

		// final_pixel_cnt = final_pixel_cnt1
		// fmt.Printf("windth -- final_pixel_cnt1: %v, final_pixel_cnt2: %v\n", final_pixel_cnt1, final_pixel_cnt2)

		if final_pixel_cnt >= 0 && (wcnt_times > short_exit || wcnt_times_long > long_exit) {
			// fmt.Println(final_pixel_wcnt)
			break
		}

	}

	cwid := math.Min(float64(width - (final_pixel_cnt - (border_px_wid * 2)) * 2), float64(width))
	return &cwid, &final_pixel_cnt
}



func crop_brd_h(img *image.Image, border_percent *float64, SimilarityThreshold *float64, short_exit_mul *float64, long_exit_mul *float64) (*float64, *int) {
	bounds := (*img).Bounds()
    width := bounds.Dx()
    height := bounds.Dy()

	short_exit := int(math.Max(float64(height) * (*short_exit_mul), 5))

	long_exit := int(math.Max(float64(height) * (*long_exit_mul), 5))

	if height < 20 {
		short_exit = 2
	}

	border_px := int(float64(height) * (*border_percent / 100))

	var final_pixel_cnt int = -1

	var cnt_times int = 0
	var cnt_times_long int = 0

	

	// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	// 	rightmostColor := (*img).At(bounds.Max.X-1, y).(color.NRGBA)
	// 	fmt.Printf("Pixel at (%d, %d) color: R=%d, G=%d, B=%d, A=%d\n", bounds.Max.X-1, y, rightmostColor.R, rightmostColor.G, rightmostColor.B, rightmostColor.A)
	// }

	
	for y := bounds.Min.Y; y < width; y++ {
		tl_col := (*img).At(0, 0).(color.NRGBA)
		// fmt.Println("h tlcol:", tl_col)
		// fmt.Println(IsSimilar(tl_col, tl_col, 10))


		cnt_times_long = 0
		for x := bounds.Min.X; x < height; x++ {
			if IsSimilar((*img).At(x-1, y).(color.NRGBA), tl_col, *SimilarityThreshold) != true {
				final_pixel_cnt = y
				cnt_times++
				cnt_times_long++
				// fmt.Println("h first", (*img).At(x-1, y).(color.NRGBA))
			} else {
				cnt_times = 0
			}
			if final_pixel_cnt >= 0 && (cnt_times > short_exit || cnt_times_long > long_exit) {
				break
			}
		}


		cnt_times_long = 0
		for x := width; x > bounds.Min.X ; x-- {
			// fmt.Println(IsSimilar((*img).At(bounds.Max.X-1, y).(color.NRGBA), tl_col, SimilarityThreshold))
			// fmt.Println(final_pixel_wcnt, x)
			if IsSimilar((*img).At(x-1, height-y-1).(color.NRGBA), tl_col, *SimilarityThreshold) != true {
				final_pixel_cnt = y
				cnt_times++
				cnt_times_long++
				// fmt.Println((*img).At(x-1, height-y-1).(color.NRGBA))
			} else {
				cnt_times = 0
			}
			if final_pixel_cnt >= 0 && (cnt_times > short_exit || cnt_times_long > long_exit) {
				break
			}
		}

		// final_pixel_cnt = int(math.Min(float64(final_pixel_cnt1), float64(final_pixel_cnt2)))
		// if final_pixel_cnt < 0 {
		// 	if final_pixel_cnt1 >= 0 || final_pixel_cnt2 >= 0 {
		// 		final_pixel_cnt = 0
		// 	}
		// }

		// // final_pixel_cnt = final_pixel_cnt1
		// fmt.Printf("height -- final_pixel_cnt1: %v, final_pixel_cnt2: %v\n", final_pixel_cnt1, final_pixel_cnt2)

		if final_pixel_cnt >= 0 && (cnt_times > short_exit || cnt_times_long > long_exit) {
			// fmt.Println(final_pixel_cnt)
			break
		}

	}

	cwid := math.Min(float64(height - (final_pixel_cnt - (border_px * 2)) * 2), float64(height))
	return &cwid, &final_pixel_cnt
}





func crop_brd(img *image.Image, border_percent *float64 , short_exit_mul *float64, long_exit_mul *float64) *image.Image {

	var SimilarityThreshold float64 = 5




	cwid, final_pixel_wcnt := crop_brd_w(img, border_percent, &SimilarityThreshold, short_exit_mul, long_exit_mul)
	chig, final_pixel_hcnt := crop_brd_h(img, border_percent, &SimilarityThreshold, short_exit_mul, long_exit_mul)


	fmt.Printf("h crop: %v, w crop: %v\n", *final_pixel_hcnt, *final_pixel_wcnt)



	if *chig <= 0 {
		*chig = 0
	}
	if *cwid <= 0 {
		*cwid = 0
	}


	croppedImg, err := cutter.Crop(*img, cutter.Config{
		Width: int(math.Round(*cwid)),
		Height: int(math.Round(*chig)),
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



func read_crop(in *string, out *string, border_p *float64 , short_exit_mul *float64, long_exit_mul *float64) {

	var img image.Image
	var err error


	in_ext := filepath.Ext(*in)
	out_ext := filepath.Ext(*out)


	file, err := os.Open(*in)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()


	decstart := time.Now()
	// Decode the WebP file
	switch in_ext {
	case ".webp":
		img, err = webp.Decode(file)
		if err != nil {
			fmt.Println("Error decoding WebP file:", err)
			return
		}
	case ".png":
		img, err = png.Decode(file)
		if err != nil {
			fmt.Println("Error decoding WebP file:", err)
			return
		}
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			fmt.Println("Error decoding WebP file:", err)
			return
		}
	default:
		if unsafe {
			img, err = webp.Decode(file)
			if err != nil {
				fmt.Println("Unknown file type and can't decode as WebP:", err)
			} else {
				break
			}
			img, err = png.Decode(file)
			if err != nil {
				fmt.Println("Unknown file type and can't decode as PNG:", err)
			} else {
				break
			}
			img, err = jpeg.Decode(file)
			if err != nil {
				fmt.Println("Unknown file type and can't decode as JPEG:", err)
				log.Fatalln("exhausted all decoding options exiting")
			} else {
				break
			}
		} else {
			log.Fatalln("can't try to decode unknown file extension with --unsafe")
		}
	}

	fmt.Println("decoding time:", time.Since(decstart))


	err = g2bwebp.Dynamic()
	if err != nil {
		fmt.Println("NON-fatal error Dynamic lib file. encoding time will be slower:\n	", err)
		// return
	}

	cropstart := time.Now()
	croppedImg := crop_brd(&img, border_p, short_exit_mul, long_exit_mul)
	fmt.Println("trim and crop time:", time.Since(cropstart))


	encstart := time.Now()
	// Create an output file
	outfile, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	switch out_ext {
	case ".webp":
		fmt.Println("webp lossless:", webp_lossless)
		err = g2bwebp.Encode(outfile, *croppedImg, g2bwebp.Options{Lossless: true, Quality: quality0_100, Method: webp_method, Exact: true})
		if err != nil {
			fmt.Println("Error encoding WebP file:", err)
			return
		}
	case ".png":
		err = png.Encode(outfile, *croppedImg)
		if err != nil {
			fmt.Println("Error encoding png file:", err)
			return
		}
	case ".jpg", ".jpeg":
		if quality0_100 < 1 {
			quality0_100 = 1
		}
		fmt.Println("jpeg quality:", quality0_100)
		err = jpeg.Encode(outfile, *croppedImg, &jpeg.Options{Quality: quality0_100})
		if err != nil {
			fmt.Println("Error encoding jpeg file:", err)
			return
		}
	default:
		log.Fatalln("can't encode unknown file extension")
	}

	fmt.Println("encoding time:", time.Since(encstart))

}


var unsafe bool = false

var webp_lossless bool = true
var webp_lossy bool
var webp_method int
var quality0_100 int

func main() {


    var input, output string
	var short_exit_mul, long_exit_mul, border_p float64

    pflag.StringVarP(&input, "input", "i", "", "file to read from")
    pflag.StringVarP(&output, "output", "o", "", "output file")
    pflag.Float64VarP(&short_exit_mul, "short_exit_mul", "s", 0.0035, "placeholder")
    pflag.Float64VarP(&long_exit_mul, "long_exit_mul", "l", 0.004, "placeholder")
    pflag.Float64VarP(&border_p, "border_percent", "b", 0.2, "a border percentage")

    pflag.BoolVar(&unsafe, "unsafe", false, "placeholder")

	pflag.BoolVar(&webp_lossless, "lossy", false, "lossy webp mode")
	pflag.IntVarP(&webp_method, "webp_method", "m", 6, "a border percentage")
	pflag.IntVarP(&quality0_100, "webp_quality", "q", 95, "lossy webp and jpeg quality, 0 to 100 for webp, 1 to 100 for jpeg")
	// pflag.IntVar(&jpeg_qual, "jpeg_quality", 95, "jpeg quality 0 to 100")

	webp_lossless = !webp_lossy

	if quality0_100 < 0 {
		quality0_100 = 0
	}
	if quality0_100 > 100 {
		quality0_100 = 100
	}

    pflag.Parse()
    viper.BindPFlags(pflag.CommandLine) // Bind pflag to viper

	read_crop(&input, &output, &border_p, &short_exit_mul, &long_exit_mul)
}