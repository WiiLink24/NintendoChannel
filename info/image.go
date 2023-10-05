package info

import (
	"NintendoChannel/constants"
	"bytes"
	_ "embed"
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
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
	if _, err := os.Stat(fmt.Sprintf("../images/%s/%s/%s.jpg", titleTypeToStr[titleType], regionToStr[region], gameID)); err == nil {
		data, err := os.ReadFile(fmt.Sprintf("../images/%s/%s/%s.jpg", titleTypeToStr[titleType], regionToStr[region], gameID))
		checkError(err)

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
		checkError(err)

		newImage := image.NewRGBA(image.Rect(0, 0, 384, 384))
		draw.BiLinear.Scale(newImage, newImage.Bounds(), img, img.Bounds(), draw.Over, nil)

		err = jpeg.Encode(buffer, newImage, nil)
		checkError(err)
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
