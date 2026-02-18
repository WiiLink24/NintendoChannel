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

	"golang.org/x/image/draw"
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

func (i *Info) WriteRatingDescriptor(buffer *bytes.Buffer, region constants.Region, RatingDescriptors []string) {
	// Cap
	maxDescriptors := min(len(RatingDescriptors), 7)
	for j := 0; j < maxDescriptors; j++ {
		s := RatingDescriptors[j]

		// Skip empty strings
		if s == "" {
			continue
		}

		// Find matching descriptor image
		var descriptorImage []byte
		switch region {
		case constants.Japan:
			descriptorImage = constants.CERODescriptors[s]
		case constants.PAL:
			descriptorImage = constants.PEGIDescriptors[s]
		case constants.NTSC:
			descriptorImage = constants.ESRBDescriptors[s]
		}

		// Set picture table entry
		i.Header.DetailedRatingPictureTable[j].PictureOffset = i.GetCurrentSize(buffer)
		buffer.Write(descriptorImage)
		i.Header.DetailedRatingPictureTable[j].PictureSize = uint32(len(descriptorImage))
	}
}
