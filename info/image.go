package info

import (
	"NintendoChannel/constants"
	"bytes"
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
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

func (i *Info) WriteCoverArt(buffer *bytes.Buffer, titleType constants.TitleType, region constants.Region, gameID string) {
	url := fmt.Sprintf("https://art.gametdb.com/%s/%s/%s/%s.png", titleTypeToStr[titleType], consoleToImageType[titleType], regionToStr[region], gameID)
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return
	}

	img, err := png.Decode(resp.Body)
	checkError(err)

	newImage := image.NewRGBA(image.Rect(0, 0, 384, 384))
	draw.BiLinear.Scale(newImage, newImage.Bounds(), img, img.Bounds(), draw.Over, nil)

	err = jpeg.Encode(buffer, newImage, nil)
	checkError(err)

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
