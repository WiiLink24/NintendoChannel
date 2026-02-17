package info

import (
	"NintendoChannel/common"
	"NintendoChannel/constants"
	"NintendoChannel/gametdb"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"strings"
	"unicode/utf16"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mitchellh/go-wordwrap"
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
	// Writing a blank one if it doesn't exist and not point to it
	// is way more efficient than writing everything individually
	TimePlayed          TimePlayed
	RecommendationTable RecommendationTable
}

var timePlayed = map[string]TimePlayed{}

func (i *Info) MakeInfo(fileID uint32, game *gametdb.Game, title, synopsis string, region constants.Region, language constants.Language, titleType constants.TitleType, ratingDescriptors []string, recommendations map[string]common.TitleRecommendation) {
	i.GetSupportedControllers(&game.Controllers)
	i.GetSupportedFeatures(game)
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

	if v, ok := timePlayed[game.ID[:4]]; ok {
		i.Header.TimesPlayedTableOffset = 6744
		i.TimePlayed = v
	}

	if v, ok := recommendations[game.ID[:4]]; ok {
		i.Header.RatingTableOffset = 6744 + 16
		i.MakeRecommendationTable(v)
	}

	temp := new(bytes.Buffer)
	imageBuffer := new(bytes.Buffer)
	i.WriteAll(temp, imageBuffer)

	i.Header.PictureOffset = i.GetCurrentSize(imageBuffer)
	i.WriteCoverArt(imageBuffer, titleType, region, game.ID)
	i.WriteRatingImage(imageBuffer, region)
	i.WriteRatingDescriptor(imageBuffer, region, ratingDescriptors)
	i.Header.Filesize = i.GetCurrentSize(imageBuffer)
	temp.Reset()

	i.WriteAll(temp, imageBuffer)
	crcTable := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum(temp.Bytes(), crcTable)
	i.Header.CRC32 = checksum
	temp.Reset()

	i.WriteAll(temp, imageBuffer)

	// Ensure write path exists
	config := common.GetConfig()
	err := os.MkdirAll(fmt.Sprintf("%s/infos/%d/%d", config.AssetsPath, region, language), 0777)
	common.CheckError(err)

	err = os.WriteFile(fmt.Sprintf("%s/infos/%d/%d/%d.info", config.AssetsPath, region, language, fileID), temp.Bytes(), 0666)
	common.CheckError(err)
}

func (i *Info) WriteAll(buffer, imageBuffer *bytes.Buffer) {
	err := binary.Write(buffer, binary.BigEndian, *i)
	common.CheckError(err)

	buffer.Write(imageBuffer.Bytes())
}

func (i *Info) GetCurrentSize(_buffer *bytes.Buffer) uint32 {
	buffer := bytes.NewBuffer(nil)
	i.WriteAll(buffer, _buffer)
	return uint32(buffer.Len())
}

func GetTimePlayed(ctx *context.Context, pool *pgxpool.Pool) {
	rows, err := pool.Query(*ctx, `SELECT game_id, COUNT(game_id), SUM(times_played), SUM(time_played) FROM time_played GROUP BY game_id`)
	common.CheckError(err)

	for rows.Next() {
		var gameID string
		var numberOfPlayers int
		var totalTimesPlayed int
		var totalTimePlayed int

		err = rows.Scan(&gameID, &numberOfPlayers, &totalTimesPlayed, &totalTimePlayed)
		common.CheckError(err)

		timePlayed[gameID] = TimePlayed{
			TotalTimePlayed:           uint32(totalTimePlayed / 60),
			TimeSpentPlayingPerPerson: uint32(totalTimePlayed / numberOfPlayers),
			TotalTimesPlayed:          uint32(totalTimesPlayed),
			TimesPlayedPerPerson:      uint32((float64(totalTimesPlayed / numberOfPlayers)) / 0.01),
		}
	}
}
