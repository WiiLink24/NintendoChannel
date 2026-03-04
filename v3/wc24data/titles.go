package wc24data

import (
	"NintendoChannel/common"
	"NintendoChannel/constants"
	"NintendoChannel/gametdb"
	"encoding/binary"
	"encoding/hex"
	"math"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/mitchellh/go-wordwrap"
)

type Company struct {
	CompanyID     uint32
	DeveloperName [31]uint16
	PublisherName [31]uint16
}

type Title struct {
	ID                         uint32
	TitleID                    [4]byte
	TitleType                  constants.TitleType
	Genre                      [3]byte
	CompanyOffset              uint32
	ReleaseYear                uint16
	ReleaseMonth               uint8
	ReleaseDay                 uint8
	RatingID                   uint8
	WithFriendsFemaleSecondRow uint8
	WithFriendsFemaleFirstRow  uint8
	WithFriendsMaleSecondRow   uint8
	WithFriendsMaleFirstRow    uint8
	WithFriendsAllSecondRow    uint8
	WithFriendsAllFirstRow     uint8
	HardcoreFemaleSecondRow    uint8
	HardcoreFemaleFirstRow     uint8
	HardcoreMaleSecondRow      uint8
	HardcoreMaleFirstRow       uint8
	HardcoreAllSecondRow       uint8
	HardcoreAllFirstRow        uint8
	GamersFemaleSecondRow      uint8
	GamersFemaleFirstRow       uint8
	GamersMaleSecondRow        uint8
	GamersMaleFirstRow         uint8
	GamersAllSecondRow         uint8
	GamersAllFirstRow          uint8
	OtherFlags                 uint8
	TitleName                  [31]uint16
	Subtitle                   [31]uint16
	ShortTitle                 [31]uint16
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
		// For some reason GameTDB has EC as 3 for some titles
		"3":    8,
		"EC":   8,
		"E":    9,
		"E10+": 10,
		"T":    11,
		"M":    12,
	},
	"PEGI": {
		"3":  8,
		"4":  8,
		"6":  9,
		"7":  9,
		"12": 10,
		"15": 11,
		"16": 11,
		"18": 12,
	},
}

func (w *WC24Data) MakeCompaniesTable() {
	w.Header.ManufacturerTableOffset = w.GetCurrentSize()

	// Only the Wii XML contains company data
	for _, company := range gametdb.WiiTDB.Companies.Companies {
		companyID, err := strconv.ParseUint(hex.EncodeToString([]byte(company.Code)), 16, 32)
		common.CheckError(err)

		name := company.Name
		if len(name) >= 31 {
			name = name[:30]
		}

		var finalDeveloperName [31]uint16
		developerName := utf16.Encode([]rune(name))
		copy(finalDeveloperName[:], developerName)

		table := Company{
			CompanyID:     uint32(companyID),
			DeveloperName: finalDeveloperName,
			PublisherName: finalDeveloperName,
		}

		w.CompanyTable = append(w.CompanyTable, table)
	}

	w.Header.NumberOfManufacturers = uint32(len(w.CompanyTable))
}

func (w *WC24Data) MakeTitleTable() {
	w.Header.TitleTableOffset = w.GetCurrentSize()

	// Wii
	w.GenerateTitleStruct(&gametdb.WiiTDB.Games, constants.Wii)
	// DS
	w.GenerateTitleStruct(&gametdb.DSTDB.Games, constants.NintendoDS)

	w.Header.NumberOfTitles = uint32(len(w.TitleTable))
}

func (w *WC24Data) GenerateTitleStruct(games *[]gametdb.Game, defaultTitleType constants.TitleType) {
	for _, game := range *games {
		if game.Locale == nil {
			// Game doesn't exist for this region?
			// Whatever the reason is, we have no metadata to use.
			continue
		}

		if game.Region == regionToGameTDB[w.region] || game.Region == "ALL" {
			titleType := defaultTitleType
			// (Sketch) The first locale will always be English from what I have observed
			title := game.Locale[0].Title
			/*
				fullTitle := game.Locale[0].Title
				synopsis := game.Locale[0].Synopsis
			*/
			for _, locale := range game.Locale {
				if locale.Language == langaugeToLocale[w.language] {
					if game.Type != "" {
						title = locale.Title
						/*
							fullTitle = locale.Title
							synopsis = locale.Synopsis
						*/
						titleType = constants.TitleTypeMap[game.Type]

						if titleType == constants.NES && defaultTitleType == constants.NintendoThreeDS {
							titleType = constants.NintendoThreeDS
						}
					}
				}
			}

			// We will not include mods or Gamecube games
			if game.Type == "CUSTOM" || game.Type == "GameCube" || game.Type == "Homebrew" || titleType == constants.ThreeDSDownload {
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
				releaseMonth = uint8(temp) - 1
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

			if len(title) >= 31 {
				title = title[:30]
			}

			if len(subtitle) >= 31 {
				subtitle = subtitle[:30]
			}

			var byteTitle [31]uint16
			tempTitle := utf16.Encode([]rune(title))
			copy(byteTitle[:], tempTitle)

			var byteSubtitle [31]uint16
			tempSubtitle := utf16.Encode([]rune(subtitle))
			copy(byteSubtitle[:], tempSubtitle)

			companyOffset, _ := w.GetCompany(&game)
			table := Title{
				ID:                         id,
				TitleID:                    titleID,
				TitleType:                  titleType,
				Genre:                      SetGenre(&game),
				CompanyOffset:              companyOffset,
				ReleaseYear:                releaseYear,
				ReleaseMonth:               releaseMonth,
				ReleaseDay:                 releaseDay,
				RatingID:                   GetRatingID(game.Rating),
				WithFriendsFemaleSecondRow: 0,
				WithFriendsFemaleFirstRow:  0,
				WithFriendsMaleSecondRow:   0,
				WithFriendsMaleFirstRow:    0,
				WithFriendsAllSecondRow:    0,
				WithFriendsAllFirstRow:     0,
				HardcoreFemaleSecondRow:    0,
				HardcoreFemaleFirstRow:     0,
				HardcoreMaleSecondRow:      0,
				HardcoreMaleFirstRow:       0,
				HardcoreAllSecondRow:       0,
				HardcoreAllFirstRow:        0,
				GamersFemaleSecondRow:      0,
				GamersFemaleFirstRow:       0,
				GamersMaleSecondRow:        0,
				GamersMaleFirstRow:         0,
				GamersAllSecondRow:         0,
				GamersAllFirstRow:          0,
				OtherFlags:                 0,
				TitleName:                  byteTitle,
				Subtitle:                   byteSubtitle,
				ShortTitle:                 [31]uint16{},
			}

			table.PopulateCriteria(w, game.ID[:4])
			table.DetermineOtherFlags(game)

			w.TitleTable = append(w.TitleTable, table)
			// TODO: Figure out how this info format works.
			/*
				i := info.Info{}
				i.MakeHeader(titleID, game.Controllers.Players, companyID, table.TitleType, table.ReleaseYear, table.ReleaseMonth, table.ReleaseDay)
				i.RatingID = table.RatingID
				i.MakeInfo(id, &game, fullTitle, synopsis, w.region, w.language, defaultTitleType)
			*/
		}
	}
}

func (w *WC24Data) GetCompany(game *gametdb.Game) (uint32, uint32) {
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
				common.CheckError(err)
				return w.Header.ManufacturerTableOffset + (128 * uint32(i)), uint32(intCompanyID)
			}
		} else {
			if strings.Contains(game.Publisher, company.Name) {
				intcompanyID, err := strconv.ParseUint(hex.EncodeToString([]byte(company.Code)), 16, 32)
				common.CheckError(err)

				return w.Header.ManufacturerTableOffset + (128 * uint32(i)), uint32(intcompanyID)
			}
		}
	}

	// If all fails, default to Nintendo
	return w.Header.ManufacturerTableOffset, 12337
}

// PopulateCriteria fills the bitfield entries in Title
func (t *Title) PopulateCriteria(l *WC24Data, gameId string) {
	if _, ok := recommendations[gameId]; !ok {
		return
	}

	// First we will go after the `All` category.
	// First 4 entries in the tables are the upper half.
	for i := 0; i < 4; i++ {
		// If it is none, then the bit will be set to nothing.
		if recommendations[gameId].AllRecommendations[i].IsGamers == constants.True {
			t.GamersAllFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsGamers == constants.False {
			t.GamersAllFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].AllRecommendations[i].IsGamers == constants.True {
			t.GamersAllSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsGamers == constants.False {
			t.GamersAllSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 0; i < 4; i++ {
		if recommendations[gameId].AllRecommendations[i].IsHardcore == constants.True {
			t.HardcoreAllFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsHardcore == constants.False {
			t.HardcoreAllFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].AllRecommendations[i].IsHardcore == constants.True {
			t.HardcoreAllSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsHardcore == constants.False {
			t.HardcoreAllSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 0; i < 4; i++ {
		if recommendations[gameId].AllRecommendations[i].IsWithFriends == constants.True {
			t.WithFriendsAllFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsWithFriends == constants.False {
			t.WithFriendsAllFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].AllRecommendations[i].IsWithFriends == constants.True {
			t.WithFriendsAllSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsWithFriends == constants.False {
			t.WithFriendsAllSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	// Next is male.
	for i := 0; i < 4; i++ {
		// If it is none, then the bit will be set to nothing.
		if recommendations[gameId].AllRecommendations[i].IsGamers == constants.True {
			t.GamersMaleFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsGamers == constants.False {
			t.GamersMaleFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].AllRecommendations[i].IsGamers == constants.True {
			t.GamersMaleSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsGamers == constants.False {
			t.GamersMaleSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 0; i < 4; i++ {
		if recommendations[gameId].AllRecommendations[i].IsHardcore == constants.True {
			t.HardcoreMaleFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsHardcore == constants.False {
			t.HardcoreMaleFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].AllRecommendations[i].IsHardcore == constants.True {
			t.HardcoreMaleSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsHardcore == constants.False {
			t.HardcoreMaleSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 0; i < 4; i++ {
		if recommendations[gameId].AllRecommendations[i].IsWithFriends == constants.True {
			t.WithFriendsMaleFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsWithFriends == constants.False {
			t.WithFriendsMaleFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].AllRecommendations[i].IsWithFriends == constants.True {
			t.WithFriendsMaleSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsWithFriends == constants.False {
			t.WithFriendsMaleSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	// Finally is female
	for i := 0; i < 4; i++ {
		// If it is none, then the bit will be set to nothing.
		if recommendations[gameId].AllRecommendations[i].IsGamers == constants.True {
			t.GamersFemaleFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsGamers == constants.False {
			t.GamersFemaleFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].AllRecommendations[i].IsGamers == constants.True {
			t.GamersFemaleSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsGamers == constants.False {
			t.GamersFemaleSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 0; i < 4; i++ {
		if recommendations[gameId].AllRecommendations[i].IsHardcore == constants.True {
			t.HardcoreFemaleFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsHardcore == constants.False {
			t.HardcoreFemaleFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].AllRecommendations[i].IsHardcore == constants.True {
			t.HardcoreFemaleSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsHardcore == constants.False {
			t.HardcoreFemaleSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 0; i < 4; i++ {
		if recommendations[gameId].AllRecommendations[i].IsWithFriends == constants.True {
			t.WithFriendsFemaleFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsWithFriends == constants.False {
			t.WithFriendsFemaleFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].AllRecommendations[i].IsWithFriends == constants.True {
			t.WithFriendsFemaleSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].AllRecommendations[i].IsWithFriends == constants.False {
			t.WithFriendsFemaleSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}
}

func (t *Title) DetermineOtherFlags(game gametdb.Game) {
	value := 0xFF
	mask := 0

	// One player or Multiplayer
	if game.Controllers.Players > 1 {
		mask |= 1
	}

	// Is the game online
	isOnline := false
	for _, s := range game.Features.Feature {
		if strings.Contains(s, "online") {
			isOnline = true
		}
	}

	if isOnline {
		mask |= 12
	}

	// TODO: We don't have titles with any videos at the moment
	t.OtherFlags = uint8(value & mask)
}

func SetGenre(game *gametdb.Game) [3]byte {
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

func (w *WC24Data) MakeNewTitleTable() {
	w.Header.NewTitleTableOffset = w.GetCurrentSize()

	titleSize := uint32(binary.Size(Title{}))
	for i, title := range w.TitleTable {
		if title.ReleaseYear != 0xFFFF && title.ReleaseMonth != 0xFF &&
			title.ReleaseDay != 0xFF && title.ReleaseYear >= 2019 {
			offset := w.Header.TitleTableOffset + uint32(i)*titleSize
			w.NewTitleTable = append(w.NewTitleTable, offset)
		}
	}

	w.Header.NumberOfNewTitles = uint32(len(w.NewTitleTable))
}

func GetRatingID(rating gametdb.Rating) uint8 {
	if rating.Value == "" {
		// Default to E/7/B
		return 9
	}

	return gameTDBRatingToRatingID[rating.Type][rating.Value]
}
