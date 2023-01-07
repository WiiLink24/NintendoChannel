package info

import (
	"NintendoChannel/constants"
	"NintendoChannel/gametdb"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/mitchellh/go-wordwrap"
	"hash/crc32"
	"log"
	"os"
	"strings"
	"unicode/utf16"
)

type Info struct {
	Header               Header
	SupportedControllers SupportedControllers
	SupportedFeatures    SupportedFeatures
	SupportedLanguages   SupportedLanguages
	_                    [10]byte
	Title                [31]uint16
	Subtitle             [31]uint16
	ShortTitle           [31]uint16
	DescriptionText      [3][41]uint16
	GenreText            [29]uint16
	PlayersText          [41]uint16
	PeripheralsText      [44]uint16
	_                    [80]byte
	DisclaimerText       [2400]uint16
	RatingID             uint8
	DistributionDateText [41]uint16
	WiiPointsText        [41]uint16
	CustomText           [10][41]uint16
}

func (i *Info) MakeInfo(fileID uint32, game *gametdb.Game, title, synopsis string, region constants.Region, language constants.Language, titleType constants.TitleType) {
	// Make other fields
	i.GetSupportedControllers(&game.Controllers)
	i.GetSupportedFeatures(&game.Features)
	i.GetSupportedLanguages(game.Languages)

	// Make title clean
	if strings.Contains(title, ": ") {
		splitTitle := strings.Split(title, ": ")
		copy(i.Title[:], utf16.Encode([]rune(splitTitle[0])))
		copy(i.Subtitle[:], utf16.Encode([]rune(splitTitle[1])))
	} else if strings.Contains(title, " - ") {
		splitTitle := strings.Split(title, " - ")
		copy(i.Title[:], utf16.Encode([]rune(splitTitle[0])))
		copy(i.Subtitle[:], utf16.Encode([]rune(splitTitle[1])))
	} else if len(title) > 31 {
		wrappedTitle := wordwrap.WrapString(title, 31)
		for p, s := range strings.Split(wrappedTitle, "\n") {
			switch p {
			case 0:
				copy(i.Title[:], utf16.Encode([]rune(s)))
				break
			case 1:
				copy(i.Subtitle[:], utf16.Encode([]rune(s)))
				break
			default:
				break
			}
		}
	} else {
		copy(i.Title[:], utf16.Encode([]rune(title)))
	}

	// Make synopsis
	wrappedSynopsis := strings.Split(wordwrap.WrapString(synopsis, 40), "\n")
	if len(wrappedSynopsis) <= 3 {
		for i2, s := range wrappedSynopsis {
			copy(i.DescriptionText[i2][:], utf16.Encode([]rune(s)))
		}
	} else {
		for i2, s := range wrappedSynopsis {
			if i2 == 10 {
				break
			} else if i2 == 9 {
				s = strings.Split(s, ".")[0] + "."
			}

			copy(i.CustomText[i2][:], utf16.Encode([]rune(s)))
		}
	}

	// Write the online players text if any
	if game.Features.OnlinePlayers != 0 {
		temp := []uint16{0, 0}
		copy(i.PlayersText[:], append(temp, utf16.Encode([]rune(fmt.Sprintf("%d Players (Online)", game.Features.OnlinePlayers)))...))
	}

	copy(i.DisclaimerText[:], utf16.Encode([]rune("Game information is provided by GameTDB.")))

	temp := new(bytes.Buffer)
	imageBuffer := new(bytes.Buffer)
	i.WriteAll(temp, imageBuffer)

	i.Header.PictureOffset = i.GetCurrentSize(imageBuffer)
	i.WriteCoverArt(imageBuffer, titleType, region, game.ID)
	i.WriteRatingImage(imageBuffer, region)
	i.Header.Filesize = i.GetCurrentSize(imageBuffer)
	temp.Reset()

	i.WriteAll(temp, imageBuffer)
	crcTable := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum(temp.Bytes(), crcTable)
	i.Header.CRC32 = checksum
	temp.Reset()

	i.WriteAll(temp, imageBuffer)
	err := os.WriteFile(fmt.Sprintf("./infos/%d/%d/%d.info", region, language, fileID), temp.Bytes(), 0666)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Nintendo Channel info file generator has encountered a fatal error! Reason: %v\n", err)
	}
}

func (i *Info) WriteAll(buffer, imageBuffer *bytes.Buffer) {
	err := binary.Write(buffer, binary.BigEndian, *i)
	checkError(err)

	buffer.Write(imageBuffer.Bytes())
}

func (i *Info) GetCurrentSize(_buffer *bytes.Buffer) uint32 {
	buffer := bytes.NewBuffer(nil)
	i.WriteAll(buffer, _buffer)
	return uint32(buffer.Len())
}
