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
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gen2brain/avif"
	"github.com/gen2brain/heic"
	"github.com/gen2brain/jpegxl"
	g2bwebp "github.com/gen2brain/webp"
	"github.com/oliamb/cutter"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/image/webp"
)



type px_range struct {
	lo_h int
	hi_h int

	lo_w int
	hi_w int
}




func in_range(x, y *int, r *px_range) bool {
	return (*y > r.lo_h && *y < r.hi_h) && (*x > r.lo_w && *x < r.hi_w)
}


var in_ranges_ir bool
var in_ranges_pxr px_range
func in_ranges(x, y *int, r *[]px_range) bool {


	for _, in_ranges_pxr = range *r {
		in_ranges_ir = in_range(x, y, &in_ranges_pxr)
		if in_ranges_ir {
			return true
		}
	}

	return false
}






func uint32_abs(n1, n2 *uint32) uint32 {
	if *n1 > *n2 {
		return *n1 - *n2
	} else {
		return *n2 - *n1
	}
}


var iss_r1, iss_g1, iss_b1, iss_a1 uint32
var iss_r2, iss_g2, iss_b2, iss_a2 uint32
var iss_flr1, iss_flg1, iss_flb1, iss_fla1 float64
var iss_flr2, iss_flg2, iss_flb2, iss_fla2 float64
// func IsSimilar(c1 color.Color, c2 *color.Color, SimilarityThreshold *float64) bool {
// 	iss_r1, iss_g1, iss_b1, iss_a1 = c1.RGBA()
// 	iss_r2, iss_g2, iss_b2, iss_a2 = (*c2).RGBA()
// 	iss_flr1 = float64(iss_r1)
// 	iss_flg1 = float64(iss_g1)
// 	iss_flb1 = float64(iss_b1)
// 	iss_fla1 = float64(iss_a1)

// 	iss_flr2 = float64(iss_r2)
// 	iss_flg2 = float64(iss_g2)
// 	iss_flb2 = float64(iss_b2)
// 	iss_fla2 = float64(iss_a2)

// 	return math.Abs(iss_flr1-iss_flr2) <= *SimilarityThreshold ||
// 		math.Abs(iss_flg1-iss_flg2) <= *SimilarityThreshold ||
// 		math.Abs(iss_flb1-iss_flb2) <= *SimilarityThreshold
// }
func IsSimilar(c1 color.Color, c2 *color.Color, SimilarityThreshold *uint32) bool {
	iss_r1, iss_g1, iss_b1, iss_a1 = c1.RGBA()
	iss_r2, iss_g2, iss_b2, iss_a2 = (*c2).RGBA()

	return uint32_abs(&iss_r1, &iss_r2) <= *SimilarityThreshold ||
		uint32_abs(&iss_g1, &iss_g2) <= *SimilarityThreshold ||
		uint32_abs(&iss_b1, &iss_b2) <= *SimilarityThreshold
}



func get_poss(l1, l2 *int, loop_type *string, width *int, height *int) (x *int, y *int) {
	switch *loop_type {
	case "r":
		retw := *width - *l1 - 1
		return &retw, l2

	case "l":
		return l1, l2

	case "t":
		retw := *l2 -1
		return &retw, l1

	case "b":
		retw := *l2 -1
		reth := *height - *l1 - 1
		return &retw, &reth
	}


	rerr := 0
	return &rerr, &rerr
}




func crop_brd_w(img *image.Image, border_percent *float64, SimilarityThreshold_fl *float64, short_exit_mul *float64, long_exit_mul *float64) (*float64, *int) {
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

	SimilarityThreshold_nonp := uint32(*SimilarityThreshold_fl)
	SimilarityThreshold := &SimilarityThreshold_nonp



	// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	// 	rightmostColor := (*img).At(bounds.Max.X-1, y).(color.NRGBA)
	// 	fmt.Printf("Pixel at (%d, %d) color: R=%d, G=%d, B=%d, A=%d\n", bounds.Max.X-1, y, rightmostColor.R, rightmostColor.G, rightmostColor.B, rightmostColor.A)
	// }

	tl_col := (*img).At(bounds.Min.X, bounds.Min.Y)
	tl_col_p := &tl_col
	for x := bounds.Min.X; x < width; x++ {
		// fmt.Println(IsSimilar(tl_col, tl_col, 10))


		wcnt_times_long = 0
		for y := bounds.Min.Y; y < height; y++ {
			if IsSimilar((*img).At(x, y), tl_col_p, SimilarityThreshold) != true {
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
			if IsSimilar((*img).At(width-x-1, y), tl_col_p, SimilarityThreshold) != true {
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



func crop_brd_h(img *image.Image, border_percent *float64, SimilarityThreshold_fl *float64, short_exit_mul *float64, long_exit_mul *float64) (*float64, *int) {
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

	SimilarityThreshold_nonp := uint32(*SimilarityThreshold_fl)
	SimilarityThreshold := &SimilarityThreshold_nonp

	// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	// 	rightmostColor := (*img).At(bounds.Max.X-1, y).(color.NRGBA)
	// 	fmt.Printf("Pixel at (%d, %d) color: R=%d, G=%d, B=%d, A=%d\n", bounds.Max.X-1, y, rightmostColor.R, rightmostColor.G, rightmostColor.B, rightmostColor.A)
	// }

	tl_col := (*img).At(bounds.Min.X, bounds.Min.Y)
	tl_col_p := &tl_col
	for y := bounds.Min.Y; y < width; y++ {

		// fmt.Println("h tlcol:", tl_col)
		// fmt.Println(IsSimilar(tl_col, tl_col, 10))


		cnt_times_long = 0
		for x := bounds.Min.X; x < height; x++ {
			if IsSimilar((*img).At(x-1, y), tl_col_p, SimilarityThreshold) != true {
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
			if IsSimilar((*img).At(x-1, height-y-1), tl_col_p, SimilarityThreshold) != true {
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




func uni_crop(img *image.Image, border_percent *float64, SimilarityThreshold_fl *float64, short_exit_mul *float64, long_exit_mul *float64, side *string, sides_map *map[string]int, ranges *[]px_range) map[string]int {
	bounds := (*img).Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	short_exit := int(math.Max(float64(width) * (*short_exit_mul), 5))

	long_exit := int(math.Max(float64(width) * (*long_exit_mul), 5))

	if width < 20 {
		short_exit = 2
	}


	var final_pixel_cnt int = -1

	var cnt_times int = 0
	var cnt_times_long int = 0

	SimilarityThreshold_nonp := uint32(*SimilarityThreshold_fl)
	SimilarityThreshold := &SimilarityThreshold_nonp



	// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	// 	rightmostColor := (*img).At(bounds.Max.X-1, y).(color.NRGBA)
	// 	fmt.Printf("Pixel at (%d, %d) color: R=%d, G=%d, B=%d, A=%d\n", bounds.Max.X-1, y, rightmostColor.R, rightmostColor.G, rightmostColor.B, rightmostColor.A)
	// }

	tl_col := (*img).At(bounds.Min.X, bounds.Min.Y)
	tl_col_p := &tl_col


	var l1, l1_max, l2, l2_max int

	switch *side {
	case "l":
		l1 = bounds.Min.X
		l1_max = width
		l2 = bounds.Min.Y
		l2_max = height

	case "r":
		l1 = bounds.Min.X
		l1_max = width
		l2 = bounds.Min.Y
		l2_max = height

	case "t":
		l1 = bounds.Min.Y
		l1_max = height
		l2 = bounds.Min.X
		l2_max = width

	case "b":
		l1 = bounds.Min.Y
		l1_max = height
		l2 = bounds.Min.X
		l2_max = width
	}

	var x, y *int

	for l1 := l1; l1 < l1_max; l1++ {
		// fmt.Println(IsSimilar(tl_col, tl_col, 10))


		cnt_times_long = 0
		for l2 := l2; l2 < l2_max; l2++ {
			x, y = get_poss(&l1, &l2, side, &width, &height)
			if in_ranges(x, y, ranges) {
				//fmt.Println("range hit", *x, *y)
				continue
			}
			// fmt.Printf("side: %v, x: %v, y: %v\n", *side, *x, *y)
			if IsSimilar((*img).At(*x, *y), tl_col_p, SimilarityThreshold) != true {
				final_pixel_cnt = l1
				cnt_times++
				cnt_times_long++
				// fmt.Println((*img).At(x, y).(color.NRGBA))
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

		// final_pixel_cnt = final_pixel_cnt1
		// fmt.Printf("windth -- final_pixel_cnt1: %v, final_pixel_cnt2: %v\n", final_pixel_cnt1, final_pixel_cnt2)

		if final_pixel_cnt >= 0 && (cnt_times > short_exit || cnt_times_long > long_exit) {
			// fmt.Println(final_pixel_wcnt)
			break
		}

	}

	(*sides_map)[*side] = final_pixel_cnt

	fmt.Println(sides_map)

	return *sides_map

}

func min_int(a int, b int) *int {
	if a < b {
		return &a
	} else {
		return &b
	}
}



func crop_brd(img *image.Image, border_percent *float64 , short_exit_mul *float64, long_exit_mul *float64, ranges *[]px_range) *image.Image {

	var SimilarityThreshold float64 = 5



	bounds := (*img).Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	border_px_wid := int(float64(width) * (*border_percent / 100))
	border_px_hi := int(float64(height) * (*border_percent / 100))












	sides := []string{"r", "l", "t", "b"}

	sides_crop := map[string]int{
		"r": 0,
		"l": 0,
		"t": 0,
		"b": 0,
	}

	var sides_crop_out map[string]int


	var final_pixel_wcnt, final_pixel_hcnt *int

	var cwid, chig *float64

	//cwid, final_pixel_wcnt := crop_brd_w(img, border_percent, &SimilarityThreshold, short_exit_mul, long_exit_mul)
	//chig, final_pixel_hcnt := crop_brd_h(img, border_percent, &SimilarityThreshold, short_exit_mul, long_exit_mul)
	for _, side := range sides {
		sides_crop_out = uni_crop(img, border_percent, &SimilarityThreshold, short_exit_mul, long_exit_mul, &side, &sides_crop, ranges)
	}


	final_pixel_wcnt = min_int(sides_crop_out["r"], sides_crop_out["l"])
	final_pixel_hcnt = min_int(sides_crop_out["t"], sides_crop_out["b"])

	cwidt := math.Min(float64(width - ((*final_pixel_wcnt) - (border_px_wid * 2)) * 2), float64(width))
	cwid = &cwidt

	chigt := math.Min(float64(height - ((*final_pixel_hcnt) - (border_px_hi * 2)) * 2), float64(height))
	chig = &chigt




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


// Convert string to image.YCbCrSubsampleRatio
func parseYCbCrSubsampleRatio(s string) (image.YCbCrSubsampleRatio, error) {
	switch s {
	case "444":
		return image.YCbCrSubsampleRatio444, nil
	case "422":
		return image.YCbCrSubsampleRatio422, nil
	case "420":
		return image.YCbCrSubsampleRatio420, nil
	case "440":
		return image.YCbCrSubsampleRatio440, nil
	case "411":
		return image.YCbCrSubsampleRatio411, nil
	case "410":
		return image.YCbCrSubsampleRatio410, nil
	default:
		return image.YCbCrSubsampleRatio(0), fmt.Errorf("unknown YCbCr subsample ratio: %s", s)
	}
}


func read_crop(in *string, out *string, border_p *float64 , short_exit_mul *float64, long_exit_mul *float64, ranges *[]px_range) {

	var img image.Image
	var err error

	mime.AddExtensionType(".webp", "image/webp")
	mime.AddExtensionType(".avif", "image/avif")
	mime.AddExtensionType(".avifs", "image/avif")
	mime.AddExtensionType(".jxl", "image/jxl")

	in_mime := mime.TypeByExtension(filepath.Ext(*in))
	out_mime := mime.TypeByExtension(filepath.Ext(*out))




	file, err := os.Open(*in)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()


	decstart := time.Now()
	// Decode the WebP file
	switch in_mime {
	case "image/webp":
		img, err = webp.Decode(file)
		if err != nil {
			fmt.Println("Error decoding WebP file:", err)
			return
		}
	case "image/avif":
		err = avif.Dynamic()
		if err != nil {
			fmt.Println("NON-fatal error Dynamic lib file. decoding time will be slower:\n	", err)
			// return
		}
		img, err = avif.Decode(file)
		if err != nil {
			fmt.Println("Error decoding WebP file:", err)
			return
		}
	case "image/jxl":
		err = jpegxl.Dynamic()
		if err != nil {
			fmt.Println("NON-fatal error Dynamic lib file. decoding time will be slower:\n	", err)
			// return
		}
		img, err = jpegxl.Decode(file)
		if err != nil {
			fmt.Println("Error decoding WebP file:", err)
			return
		}
	case "image/heif", "image/heif-sequence", "image/heic", "image/heic-sequence":
		err = heic.Dynamic()
		if err != nil {
			fmt.Println("NON-fatal error Dynamic lib file. decoding time will be slower:\n	", err)
			// return
		}
		img, err = heic.Decode(file)
		if err != nil {
			fmt.Println("Error decoding WebP file:", err)
			return
		}


	case "image/png":
		img, err = png.Decode(file)
		if err != nil {
			fmt.Println("Error decoding WebP file:", err)
			return
		}
	case "image/jpeg":
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
			log.Fatalln("can't try to decode unknown file extension without --unsafe")
		}
	}

	fmt.Println("decoding time:", time.Since(decstart))




	cropstart := time.Now()
	croppedImg := crop_brd(&img, border_p, short_exit_mul, long_exit_mul, ranges)
	fmt.Println("trim and crop time:", time.Since(cropstart))


	encstart := time.Now()
	// Create an output file
	outfile, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()


	switch out_mime {
	case "image/webp":
		err = g2bwebp.Dynamic()
		if err != nil {
			fmt.Println("NON-fatal error Dynamic lib file. encoding time will be slower:\n	", err)
			// return
		}
		fmt.Println("webp lossless:", webp_lossless)
		err = g2bwebp.Encode(outfile, *croppedImg, g2bwebp.Options{Lossless: true, Quality: quality0_100, Method: webp_method, Exact: true})
		if err != nil {
			fmt.Println("Error encoding WebP file:", err)
			return
		}
	case "image/avif":
		err = avif.Dynamic()
		if err != nil {
			fmt.Println("NON-fatal error Dynamic lib file. encoding time will be slower:\n	", err)
			// return
		}
		fmt.Printf("Quality: %v, QualityAlpha: %v, Speed: %v, ChromaSubsampling: %v\n", quality0_100, quality0_100_alpha, avif_speed, chroma_sub)
		err = avif.Encode(outfile, *croppedImg, avif.Options{Quality: quality0_100, QualityAlpha: quality0_100_alpha, Speed: avif_speed, ChromaSubsampling: chroma_sub})
		if err != nil {
			fmt.Println("Error encoding WebP file:", err)
			return
		}
	case "image/jxl":
		err = jpegxl.Dynamic()
		if err != nil {
			fmt.Println("NON-fatal error Dynamic lib file. encoding time will be slower:\n	", err)
			// return
		}

		err = jpegxl.Encode(outfile, *croppedImg, jpegxl.Options{Quality: quality0_100, Effort: jpegxl_effort})
		if err != nil {
			fmt.Println("Error encoding jxl file:", err)
			return
		}
	case "image/heif", "image/heif-sequence", "image/heic", "image/heic-sequence":
		log.Fatalln("can't encode to heic")
	case "image/png":
		err = png.Encode(outfile, *croppedImg)
		if err != nil {
			fmt.Println("Error encoding png file:", err)
			return
		}
	case "image/jpeg":
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





func parse_excludes(s string) []px_range {
	var fin_range []px_range

	s = strings.ReplaceAll(s, " ", "")
	ranges := strings.Split(s, ";")


	for _, rang := range ranges {
		rang_spl := strings.Split(rang, "-")
		r1 := strings.Split(rang_spl[0], ",")
		r2 := strings.Split(rang_spl[1], ",")

		x1, err := strconv.Atoi(r1[0])
		if err != nil {
			log.Fatalln(err)
		}

		y1, err := strconv.Atoi(r1[1])
		if err != nil {
			log.Fatalln(err)
		}



		x2, err := strconv.Atoi(r2[0])
		if err != nil {
			log.Fatalln(err)
		}

		y2, err := strconv.Atoi(r2[1])
		if err != nil {
			log.Fatalln(err)
		}


		var low_x int
		var hi_x int

		var low_y int
		var hi_y int

		if x1 < x2 {
			low_x = x1
			hi_x = x2
		} else {
			low_x = x2
			hi_x = x1
		}

		if y1 < y2 {
			low_y = y1
			hi_y = y2
		} else {
			low_y = y2
			hi_y = y1
		}



		fin_range = append(fin_range, px_range{lo_w: low_x, lo_h: low_y, hi_w: hi_x, hi_h: hi_y})
	}

	return fin_range
}







var unsafe bool = false

var webp_lossless bool = true
var webp_lossy bool
var webp_method int
var avif_speed int
var jpegxl_effort int
var chroma_sub_str string
var chroma_sub image.YCbCrSubsampleRatio
var quality0_100 int
var quality0_100_alpha int

func main() {
	var err error

	var input, output string
	var short_exit_mul, long_exit_mul, border_p float64


	var input_ex_ranges string
	var ex_ranges []px_range

	pflag.StringVarP(&input, "input", "i", "", "file to read from")
	pflag.StringVarP(&output, "output", "o", "", "output file")
	pflag.Float64VarP(&short_exit_mul, "short_exit_mul", "s", 0.003, "placeholder")
	pflag.Float64VarP(&long_exit_mul, "long_exit_mul", "l", 0.004, "placeholder")
	pflag.Float64VarP(&border_p, "border_percent", "b", 0.2, "a border percentage")

	pflag.BoolVar(&unsafe, "unsafe", false, "placeholder")

	pflag.BoolVar(&webp_lossy, "lossy", false, "lossy webp mode")
	pflag.IntVarP(&webp_method, "webp_method", "m", 6, "webp compression method (0=fastest, 6=slowest)")
	pflag.IntVar(&avif_speed, "avif_speed", 0, "Speed in the range [0,10]. Slower should make for a better quality image in less bytes. lower is slower")
	pflag.IntVar(&jpegxl_effort, "jpegxl_effort", 7, "Effort in the range [1,10]. Sets encoder effort/speed level without affecting decoding speed. Default is 7.")
	pflag.StringVar(&chroma_sub_str, "chroma_sub", "444", "Chroma subsampling, 444|422|420. applys to avif")
	pflag.StringVarP(&input_ex_ranges, "exclude_ranges", "e", "", "px ranges to exclude")
	pflag.IntVarP(&quality0_100, "quality", "q", 100, "lossy webp and jpeg quality, 0 to 100 for webp, avif, jpeg xl, heic. 1 to 100 for jpeg\nQuality of 100 implies lossless for webp, jpeg xl, and avif")
	pflag.IntVarP(&quality0_100_alpha, "quality_alpha", "a", 100, "alpha quality. avif,")
	// pflag.IntVar(&jpeg_qual, "jpeg_quality", 95, "jpeg quality 0 to 100")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine) // Bind pflag to viper




	fmt.Printf("\n\n========================================\n========================================\nmaking file: %v\n\n\n", output)





	if input_ex_ranges != "" {
		ex_ranges = parse_excludes(input_ex_ranges)
		fmt.Println("exclude:", ex_ranges)
	}


	fmt.Println(webp_lossy)
	webp_lossless = !webp_lossy

	chroma_sub, err = parseYCbCrSubsampleRatio(chroma_sub_str)
	if err != nil {
		fmt.Println(err)
		fmt.Println("setting chrome subsampling to 422")
		chroma_sub = image.YCbCrSubsampleRatio422
	}

	if quality0_100 < 0 {
		quality0_100 = 0
	}
	if quality0_100 > 100 {
		quality0_100 = 100
	}



	read_crop(&input, &output, &border_p, &short_exit_mul, &long_exit_mul, &ex_ranges)



	fmt.Print("\n\n========================================\n========================================\n")
}
