package main

import (
	"NintendoChannel/common"
	"NintendoChannel/constants"
	"NintendoChannel/gametdb"
	"bytes"
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"testing"
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

var regionToGameTDB = map[constants.Region]string{
	constants.NTSC:  "NTSC-U",
	constants.PAL:   "PAL",
	constants.Japan: "NTSC-J",
}

func TestGetAllImages(t *testing.T) {
	config := common.GetConfig()

	for _, region := range constants.Regions {
		gametdb.PrepareGameTDB(config)
		for i, games := range [][]gametdb.Game{gametdb.DSTDB.Games, gametdb.WiiTDB.Games, gametdb.ThreeDSTDB.Games} {
			for _, game := range games {
				titleType := constants.NintendoDS
				if i == 1 {
					titleType = constants.Wii
				} else if i == 2 {
					titleType = constants.NintendoThreeDS
				}

				if game.Type == "CUSTOM" || game.Type == "GameCube" || game.Type == "Homebrew" {
					continue
				}

				if game.Region != regionToGameTDB[region.Region] {
					continue
				}

				if _, err := os.Stat(fmt.Sprintf("%s/%s/%s/%s.jpg", config.ImagesPath, titleTypeToStr[titleType], regionToStr[region.Region], game.ID)); err == nil {
					fmt.Println("Skipping ", game.ID)
					continue
				}

				url := fmt.Sprintf("https://art.gametdb.com/%s/%s/%s/%s.png", titleTypeToStr[titleType], consoleToImageType[titleType], regionToStr[region.Region], game.ID)
				resp, err := http.Get(url)
				if err != nil {
					panic(err)
				}

				if resp.StatusCode != http.StatusOK {
					continue
				}

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
				draw.Draw(newImage, img.Bounds().Add(offset), img, image.Point{}, draw.Src)

				buffer := new(bytes.Buffer)
				err = jpeg.Encode(buffer, newImage, nil)
				common.CheckError(err)

				err = os.MkdirAll(fmt.Sprintf("%s/%s/%s", config.ImagesPath, titleTypeToStr[titleType], regionToStr[region.Region]), 0777)
				if err != nil {
					t.Error(err)
				}

				err = os.WriteFile(fmt.Sprintf("%s/%s/%s/%s.jpg", config.ImagesPath, titleTypeToStr[titleType], regionToStr[region.Region], game.ID), buffer.Bytes(), 0666)
				if err != nil {
					t.Error(err)
				}

				resp.Body.Close()
			}
		}
	}
}

func resize(origImage image.Image, x, y int) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, x, y))
	draw.BiLinear.Scale(newImage, newImage.Bounds(), origImage, origImage.Bounds(), draw.Over, nil)
	return newImage
}
