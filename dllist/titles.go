package dllist

import (
	"NintendoChannel/constants"
	"NintendoChannel/gametdb"
	"NintendoChannel/info"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/mitchellh/go-wordwrap"
	"strconv"
	"strings"
	"unicode/utf16"
)

// CompanyTable represents a company in the dllist.bin
type CompanyTable struct {
	CompanyID     uint32
	DeveloperName [31]uint16
	PublisherName [31]uint16
}

type TitleTable struct {
	ID        uint32
	TitleID   [4]byte
	TitleType constants.TitleType
	// TODO: Implement genres
	Genre         [3]byte
	CompanyOffset uint32
	ReleaseYear   uint16
	ReleaseMonth  uint8
	ReleaseDay    uint8
	RatingID      uint8
	Unknown       [2]byte
	// TODO: Giant bitfield, extremely low priority however implement later.
	HardcoreBitField uint32
	FriendsBitField  uint32
	Unknown3BitField uint32
	Unknown4         uint16
	Unknown5         uint16
	Unknown6         uint8
	Unknown7         uint32
	Unknown8         uint32
	MedalType        uint8
	Unknown9         uint8
	TitleName        [31]uint16
	Subtitle         [31]uint16
	ShortTitle       [31]uint16
}

func (l *List) MakeCompaniesTable() {
	l.Header.CompanyTableOffset = l.GetCurrentSize()

	// Only the Wii XML contains company data
	for _, company := range gametdb.WiiTDB.Companies.Companies {
		companyID, err := strconv.ParseUint(hex.EncodeToString([]byte(company.Code)), 16, 32)
		checkError(err)

		var finalDeveloperName [31]uint16
		developerName := utf16.Encode([]rune(company.Name))
		copy(finalDeveloperName[:], developerName)

		table := CompanyTable{
			CompanyID:     uint32(companyID),
			DeveloperName: finalDeveloperName,
			PublisherName: finalDeveloperName,
		}

		l.CompaniesTable = append(l.CompaniesTable, table)
	}

	l.Header.NumberOfCompanyTables = uint32(len(l.CompaniesTable))
}

var langaugeToLocale = map[constants.Language]string{
	constants.Japanese: "JA",
	constants.English:  "EN",
	constants.German:   "DE",
	constants.French:   "FR",
	constants.Spanish:  "ES",
	constants.Italian:  "IT",
	constants.Dutch:    "NL",
}

var regionToGameTDB = map[constants.Region]string{
	constants.NTSC:  "NTSC-U",
	constants.PAL:   "PAL",
	constants.Japan: "NTSC-J",
}

var gameTDBRatingToRatingID = map[string]map[string]uint8{
	"CERO": {
		"A": 8,
		"B": 9,
		"C": 10,
		"D": 11,
		"Z": 12,
	},
	"ESRB": {
		"EC":   8,
		"E":    9,
		"E10+": 10,
		"T":    11,
		"M":    12,
	},
	"PEGI": {
		"3":  8,
		"7":  9,
		"12": 10,
		"16": 11,
		"18": 12,
	},
}

func (l *List) MakeTitleTable() {
	l.Header.TitleTableOffset = l.GetCurrentSize()

	// Wii
	l.GenerateTitleStruct(&gametdb.WiiTDB.Games, constants.Wii)
	// DS
	l.GenerateTitleStruct(&gametdb.DSTDB.Games, constants.NintendoDS)
	// 3DS
	l.GenerateTitleStruct(&gametdb.ThreeDSTDB.Games, constants.NintendoThreeDS)

	l.Header.NumberOfTitleTables = uint32(len(l.TitleTable))
}

func (l *List) GenerateTitleStruct(games *[]gametdb.Game, defaultTitleType constants.TitleType) {
	for _, game := range *games {
		if game.Region == regionToGameTDB[l.region] || game.Region == "ALL" {
			titleType := defaultTitleType
			// (Sketch) The first locale will always be English from what I have observed
			title := game.Locale[0].Title
			fullTitle := game.Locale[0].Title
			synopsis := game.Locale[0].Synopsis
			for _, locale := range game.Locale {
				if locale.Language == langaugeToLocale[l.language] {
					if game.Type != "" {
						title = locale.Title
						fullTitle = locale.Title
						synopsis = locale.Synopsis
						titleType = constants.TitleTypeMap[game.Type]
					}
				}
			}

			// We will not include mods or Gamecube games
			if game.Type == "CUSTOM" || game.Type == "GameCube" {
				continue
			}

			var titleID [4]byte
			copy(titleID[:], game.ID)

			// Wii, DS and 3DS games may share the same IDs are one another. XOR to avoid conflict.
			id := binary.BigEndian.Uint32(titleID[:])
			if defaultTitleType == constants.NintendoDS {
				id ^= 0x22222222
			} else if defaultTitleType == constants.NintendoThreeDS {
				id ^= 0x33333333
			}

			var releaseYear uint16 = 0xFFFF
			if game.ReleaseDate.Year != "" {
				temp, _ := strconv.ParseUint(game.ReleaseDate.Year, 10, 32)
				releaseYear = uint16(temp)
			}

			var releaseMonth uint8 = 0xFF
			if game.ReleaseDate.Month != "" {
				temp, _ := strconv.ParseUint(game.ReleaseDate.Month, 10, 32)
				releaseMonth = uint8(temp)
			}

			var releaseDay uint8 = 0xFF
			if game.ReleaseDate.Day != "" {
				temp, _ := strconv.ParseUint(game.ReleaseDate.Day, 10, 32)
				releaseDay = uint8(temp)
			}

			subtitle := ""
			if len(title) > 30 {
				wrappedTitle := wordwrap.WrapString(title, 30)
				for i, s := range strings.Split(wrappedTitle, "\n") {
					switch i {
					case 0:
						title = s
						break
					case 1:
						subtitle = s
						break
					default:
						break
					}
				}
			}

			var byteTitle [31]uint16
			tempTitle := utf16.Encode([]rune(title))
			copy(byteTitle[:], tempTitle)

			var byteSubtitle [31]uint16
			tempSubtitle := utf16.Encode([]rune(subtitle))
			copy(byteSubtitle[:], tempSubtitle)

			companyOffset, companyID := l.GetCompany(&game)
			table := TitleTable{
				ID:               id,
				TitleID:          titleID,
				TitleType:        titleType,
				Genre:            l.SetGenre(&game),
				CompanyOffset:    companyOffset,
				ReleaseYear:      releaseYear,
				ReleaseMonth:     releaseMonth,
				ReleaseDay:       releaseDay,
				RatingID:         GetRatingID(game.Rating),
				Unknown:          [2]byte{0x8, 0x20},
				HardcoreBitField: 536872960,
				FriendsBitField:  2863311530,
				Unknown3BitField: 2863311530,
				Unknown4:         43690,
				Unknown5:         170,
				Unknown6:         168,
				Unknown7:         50331648,
				Unknown8:         0,
				MedalType:        0,
				Unknown9:         222,
				TitleName:        byteTitle,
				Subtitle:         byteSubtitle,
				ShortTitle:       [31]uint16{},
			}

			l.TitleTable = append(l.TitleTable, table)

			// Write all our static data first
			i := info.Info{}
			i.MakeHeader(titleID, game.Controllers.Players, companyID, table.TitleType, table.ReleaseYear, table.ReleaseMonth, table.ReleaseDay)
			i.RatingID = table.RatingID
			i.MakeInfo(&game, fullTitle, synopsis, l.region, defaultTitleType)
		}
	}
}

func GetRatingID(rating gametdb.Rating) uint8 {
	if rating.Value == "" {
		// Default to E/7/B
		return 9
	}

	return gameTDBRatingToRatingID[rating.Type][rating.Value]
}

func (l *List) GetCompany(game *gametdb.Game) (uint32, uint32) {
	isDiscGame := false
	companyID := ""
	// This first method of retrieving the company is the most accurate. However, it only works with disc games.
	if len(game.ID) != 4 {
		isDiscGame = true
		companyID = game.ID[4:]
	}

	for i, company := range gametdb.WiiTDB.Companies.Companies {
		if isDiscGame {
			if companyID == company.Code {
				intCompanyID, err := strconv.ParseUint(hex.EncodeToString([]byte(company.Code)), 16, 32)
				checkError(err)
				return l.Header.CompanyTableOffset + (128 * uint32(i)), uint32(intCompanyID)
			}
		} else {
			if strings.Contains(game.Publisher, company.Name) {
				intcompanyID, err := strconv.ParseUint(hex.EncodeToString([]byte(company.Code)), 16, 32)
				checkError(err)

				return l.Header.CompanyTableOffset + (128 * uint32(i)), uint32(intcompanyID)
			}
		}
	}

	// If all fails, default to Nintendo
	return l.Header.CompanyTableOffset, 12337
}

func (l *List) SetGenre(game *gametdb.Game) [3]byte {
	gameTDBToGenre := map[string]uint8{
		"arcade":               15,
		"party":                13,
		"puzzle":               5,
		"action":               1,
		"2D platformer":        1,
		"3D platformer":        1,
		"shooter":              12,
		"first-person shooter": 12,
		"third-person shooter": 12,
		"rail shooter":         12,
		"run and gun":          12,
		"shoot 'em up":         12,
		"stealth action":       1,
		"survival horror":      1,
		"sports":               4,
		"adventure":            2,
		"hidden object":        13,
		"interactive fiction":  2,
		"interactive movie":    2,
		"point-and-click":      13,
		"music":                10,
		"rhythm":               10,
		"dance":                10,
		"karaoke":              10,
		"racing":               7,
		"fighting":             14,
		"simulation":           9,
		"role-playing":         6,
		"strategy":             8,
		"traditional":          11,
		"health":               3,
		"others":               13,
	}

	genre := [3]byte{0, 0, 0}
	for i, s := range strings.Split(game.Genre, ",") {
		if i == 3 {
			break
		}

		genre[i] = gameTDBToGenre[s]
	}

	return genre
}

func (l *List) MakeNewTitleTable() {
	// TODO: Figure out a way to get the newest titles
	l.Header.NewTitleTableOffset = l.GetCurrentSize()
	l.NewTitleTable = append(l.NewTitleTable, l.Header.TitleTableOffset)
	l.Header.NumberOfNewTitleTables = 1
}
