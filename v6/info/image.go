package info

import (
	"NintendoChannel/common"
	"NintendoChannel/constants"
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strings"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var regionToStr = map[constants.Region]string{
	constants.Japan: "JA",
	constants.PAL:   "EN",
	constants.NTSC:  "US",
}
var titleTypeToStr = map[constants.TitleType]string{
	constants.Wii:             "wii",
	constants.NintendoDS:      "ds",
	constants.NintendoThreeDS: "3ds",
}

var consoleToImageType = map[constants.TitleType]string{
	constants.Wii:             "cover",
	constants.NintendoDS:      "box",
	constants.NintendoThreeDS: "box",
}

var consoleToTempImageType = map[constants.TitleType][]byte{
	constants.Wii:             PlaceholderWii,
	constants.NintendoDS:      PlaceholderDS,
	constants.NintendoThreeDS: Placeholder3DS,
}

//go:embed 3ds.jpg
var Placeholder3DS []byte

//go:embed ds.jpg
var PlaceholderDS []byte

//go:embed wii.jpg
var PlaceholderWii []byte

func (i *Info) WriteCoverArt(buffer *bytes.Buffer, titleType constants.TitleType, region constants.Region, gameID string) {
	// Check if it exists on disk first.
	path := fmt.Sprintf("%s/%s/%s/%s.jpg", common.GetConfig().ImagesPath, titleTypeToStr[titleType], regionToStr[region], gameID)
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		common.CheckError(err)

		buffer.Write(data)
		i.Header.PictureSize = uint32(buffer.Len())
		return
	}

	url := fmt.Sprintf("https://art.gametdb.com/%s/%s/%s/%s.png", titleTypeToStr[titleType], consoleToImageType[titleType], regionToStr[region], gameID)
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		buffer.Write(consoleToTempImageType[titleType])
	} else {
		img, err := png.Decode(resp.Body)
		common.CheckError(err)

		// Some resizing on the image to make it not look as stretched
		x, y := img.Bounds().Dx(), img.Bounds().Dy()

		if titleType != constants.NintendoThreeDS && titleType != constants.NintendoDS {
			img = resize(img, int(float64(x)*(384.0/float64(y))), 384)
		} else {
			img = resize(img, 384, int(float64(y)*(384.0/float64(x))))
		}

		offsetX := (384 - img.Bounds().Dx()) / 2
		offsetY := (384 - img.Bounds().Dy()) / 2
		offset := image.Pt(offsetX, offsetY)

		// Creates a blank white image which will then be layered by the cover
		newImage := image.NewRGBA(image.Rect(0, 0, 384, 384))
		draw.Draw(newImage, newImage.Bounds(), &image.Uniform{C: color.RGBA{R: 255, G: 255, B: 255, A: 255}}, image.Point{}, draw.Src)
		draw.Draw(newImage, img.Bounds().Add(offset), img, image.Point{}, draw.Over)

		err = jpeg.Encode(buffer, newImage, nil)
		common.CheckError(err)

		// Cache image for future generations.
		err = os.WriteFile(path, buffer.Bytes(), 0666)
		common.CheckError(err)
	}

	i.Header.PictureSize = uint32(buffer.Len())
}

func (i *Info) WriteRatingImage(buffer *bytes.Buffer, region constants.Region) {
	i.Header.RatingPictureOffset = i.GetCurrentSize(buffer)

	regionToRatingGroup := map[constants.Region]constants.RatingGroup{
		constants.Japan: constants.CERO,
		constants.NTSC:  constants.ESRB,
		constants.PAL:   constants.PEGI,
	}

	buffer.Write(constants.Images[regionToRatingGroup[region]][i.RatingID-8])
	i.Header.RatingPictureSize = uint32(len(constants.Images[regionToRatingGroup[region]][i.RatingID-8]))
}

func resize(origImage image.Image, x, y int) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, x, y))
	draw.BiLinear.Scale(newImage, newImage.Bounds(), origImage, origImage.Bounds(), draw.Over, nil)
	return newImage
}

func (i *Info) WriteDetailedRatingImage(buffer *bytes.Buffer, region constants.Region, DetailedRatingPictureTable [7]string) {
	if region == 2 { // NTSC-U, ESRB
		// Parse the embedded ESRB font (Rodin NTLG)
		f, err := opentype.Parse(constants.ESRBRatingDescriptorFont)
		common.CheckError(err)

		// Create font face with appropriate size for ESRB content descriptors
		face, err := opentype.NewFace(f, &opentype.FaceOptions{
			Size:    15,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		common.CheckError(err)

		for j, s := range DetailedRatingPictureTable {

			// Skip empty strings to avoid creating unnecessary images
			if s == "" {
				continue
			}

			// Create 350x16 image with white background
			img := image.NewRGBA(image.Rect(0, 0, 350, 16))
			draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

			// Create the font drawer
			d := &font.Drawer{
				Dst:  img,
				Src:  image.NewUniform(color.Black),
				Face: face,
				Dot:  fixed.Point26_6{X: fixed.I(2), Y: fixed.I(13)},
			}

			// Draw the text
			words := strings.Fields(s)
			for i, word := range words {
				if len(word) > 0 {
					words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
				}
			}
			d.DrawString(strings.Join(words, " "))

			// Encode image
			var imgBuffer bytes.Buffer
			err := jpeg.Encode(&imgBuffer, img, &jpeg.Options{Quality: 100})
			common.CheckError(err)

			// Set picture table entry
			i.Header.DetailedRatingPictureTable[j].PictureOffset = i.GetCurrentSize(buffer)
			buffer.Write(imgBuffer.Bytes())
			i.Header.DetailedRatingPictureTable[j].PictureSize = uint32(imgBuffer.Len())

		}
	} else if region == 1 || region == 0 { // PAL (PEGI) and NTSC-J (CERO)

		for j, s := range DetailedRatingPictureTable {
			// Skip empty strings
			if s == "" {
				continue
			}

			// Find matching descriptor image
			var descriptorImage []byte
			var exists bool
			if region == 1 {
				descriptorImage, exists = constants.PEGIDescriptors[s]
			} else {
				descriptorImage, exists = constants.CERODescriptors[s]
			}

			// Skip if no matching descriptor found
			if !exists {
				continue
			}

			// Set picture table entry
			i.Header.DetailedRatingPictureTable[j].PictureOffset = i.GetCurrentSize(buffer)
			buffer.Write(descriptorImage)
			i.Header.DetailedRatingPictureTable[j].PictureSize = uint32(len(descriptorImage))

		}
	}
}
