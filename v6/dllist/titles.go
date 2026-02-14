package dllist

import (
	"NintendoChannel/common"
	"NintendoChannel/constants"
	"NintendoChannel/gametdb"
	"NintendoChannel/v6/info"
	"encoding/hex"
	"hash/crc32"
	"math"
	"slices"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/mitchellh/go-wordwrap"
)

// CompanyTable represents a company in the dllist.bin
type CompanyTable struct {
	CompanyID     uint32
	DeveloperName [31]uint16
	PublisherName [31]uint16
}

type TitleTable struct {
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
	Unknown7                   [4]byte
	Unknown8                   uint32
	MedalType                  constants.Medal
	Unknown9                   uint8
	TitleName                  [31]uint16
	Subtitle                   [31]uint16
	ShortTitle                 [31]uint16
}

func (l *List) MakeCompaniesTable() {
	l.Header.CompanyTableOffset = l.GetCurrentSize()

	// Only the Wii XML contains company data
	for _, company := range gametdb.WiiTDB.Companies.Companies {
		companyID, err := strconv.ParseUint(hex.EncodeToString([]byte(company.Code)), 16, 32)
		common.CheckError(err)

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

var regionToCodeTDB = map[constants.Region]byte{
	constants.NTSC:  'E',
	constants.PAL:   'P',
	constants.Japan: 'J',
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
	// Internally only created once, but it doesn't hurt to have it on the stack.
	// Required for generating a unique checksum for title ids.
	crcTable := crc32.MakeTable(crc32.IEEE)

	for _, game := range *games {
		if game.Locale == nil {
			// Game doesn't exist for this region?
			// Whatever the reason is, we have no metadata to use.
			continue
		}

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

						if titleType == constants.NES && defaultTitleType == constants.NintendoThreeDS {
							titleType = constants.NintendoThreeDS
						}
					}
				}
			}

			// We will not include mods, GameCube games, demos, DS Download Stations, or Pokemon distributions
			if game.Type == "CUSTOM" || game.Type == "GameCube" || game.Type == "Homebrew" || titleType == constants.ThreeDSDownload ||
				strings.Contains(title, "(Demo)") || strings.Contains(title, "Download") ||
				strings.Contains(title, "Distribution") || strings.Contains(title, "DSi XL") ||
				strings.Contains(title, "Exclusive") || strings.Contains(title, "Toys R Us") ||
				strings.Contains(title, "GameStop") || strings.Contains(title, "Target") ||
				strings.Contains(title, "Best Buy") || strings.Contains(title, "Walmart") ||
				strings.Contains(title, "Limited Edition") || strings.Contains(title, "Collector's Edition") ||
				strings.Contains(title, "(Beta)") || strings.Contains(title, "Relay") ||
				slices.Contains(constants.DevAppIDs, game.ID[:4]) {
				continue
			}

			if game.ID[3] != regionToCodeTDB[l.region] {
				continue
			}

			var titleID [4]byte
			copy(titleID[:], game.ID)

			// Wii, DS and 3DS games may share the same IDs are one another. XOR to avoid conflict.
			id := crc32.Checksum([]byte(game.ID), crcTable)
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

			var byteTitle [31]uint16
			tempTitle := utf16.Encode([]rune(title))
			copy(byteTitle[:], tempTitle)

			var byteSubtitle [31]uint16
			tempSubtitle := utf16.Encode([]rune(subtitle))
			copy(byteSubtitle[:], tempSubtitle)

			medal := constants.None
			if num, ok := recommendations[game.ID[:4]]; ok {
				medal = GetMedal(num.NumberOfRecommendations)
			}

			companyOffset, companyID := l.GetCompany(&game)
			table := TitleTable{
				ID:                         id,
				TitleID:                    titleID,
				TitleType:                  titleType,
				Genre:                      l.SetGenre(&game),
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
				Unknown7:                   [4]byte{0, 0, 0},
				Unknown8:                   0,
				MedalType:                  medal,
				Unknown9:                   222,
				TitleName:                  byteTitle,
				Subtitle:                   byteSubtitle,
				ShortTitle:                 [31]uint16{},
			}

			table.PopulateCriteria(l, game.ID[:4])
			table.DetermineOtherFlags(game)

			l.TitleTable = append(l.TitleTable, table)
			if !generateTitles {
				continue
			}

			i := info.Info{}
			i.MakeHeader(titleID, game.Controllers.Players, companyID, table.TitleType, table.ReleaseYear, table.ReleaseMonth, table.ReleaseDay)
			i.RatingID = table.RatingID
			i.MakeInfo(id, &game, fullTitle, synopsis, l.region, l.language, defaultTitleType, recommendations)
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
				common.CheckError(err)
				return l.Header.CompanyTableOffset + (128 * uint32(i)), uint32(intCompanyID)
			}
		} else {
			if strings.Contains(game.Publisher, company.Name) {
				intcompanyID, err := strconv.ParseUint(hex.EncodeToString([]byte(company.Code)), 16, 32)
				common.CheckError(err)

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

func GetMedal(numberOfTimesVotes int) constants.Medal {
	if numberOfTimesVotes >= 50 {
		return constants.Platinum
	} else if numberOfTimesVotes >= 35 {
		return constants.Gold
	} else if numberOfTimesVotes >= 20 {
		return constants.Silver
	} else if numberOfTimesVotes >= 15 {
		return constants.Bronze
	}

	return constants.None
}

// PopulateCriteria fills the bitfield entries in TitleTable
func (t *TitleTable) PopulateCriteria(l *List, gameId string) {
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
		if recommendations[gameId].MaleRecommendations[i].IsGamers == constants.True {
			t.GamersMaleFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].MaleRecommendations[i].IsGamers == constants.False {
			t.GamersMaleFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].MaleRecommendations[i].IsGamers == constants.True {
			t.GamersMaleSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].MaleRecommendations[i].IsGamers == constants.False {
			t.GamersMaleSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 0; i < 4; i++ {
		if recommendations[gameId].MaleRecommendations[i].IsHardcore == constants.True {
			t.HardcoreMaleFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].MaleRecommendations[i].IsHardcore == constants.False {
			t.HardcoreMaleFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].MaleRecommendations[i].IsHardcore == constants.True {
			t.HardcoreMaleSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].MaleRecommendations[i].IsHardcore == constants.False {
			t.HardcoreMaleSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 0; i < 4; i++ {
		if recommendations[gameId].MaleRecommendations[i].IsWithFriends == constants.True {
			t.WithFriendsMaleFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].MaleRecommendations[i].IsWithFriends == constants.False {
			t.WithFriendsMaleFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].MaleRecommendations[i].IsWithFriends == constants.True {
			t.WithFriendsMaleSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].MaleRecommendations[i].IsWithFriends == constants.False {
			t.WithFriendsMaleSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	// Finally is female
	for i := 0; i < 4; i++ {
		// If it is none, then the bit will be set to nothing.
		if recommendations[gameId].FemaleRecommendations[i].IsGamers == constants.True {
			t.GamersFemaleFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].FemaleRecommendations[i].IsGamers == constants.False {
			t.GamersFemaleFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].FemaleRecommendations[i].IsGamers == constants.True {
			t.GamersFemaleSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].FemaleRecommendations[i].IsGamers == constants.False {
			t.GamersFemaleSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 0; i < 4; i++ {
		if recommendations[gameId].FemaleRecommendations[i].IsHardcore == constants.True {
			t.HardcoreFemaleFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].FemaleRecommendations[i].IsHardcore == constants.False {
			t.HardcoreFemaleFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].FemaleRecommendations[i].IsHardcore == constants.True {
			t.HardcoreFemaleSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].FemaleRecommendations[i].IsHardcore == constants.False {
			t.HardcoreFemaleSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 0; i < 4; i++ {
		if recommendations[gameId].FemaleRecommendations[i].IsWithFriends == constants.True {
			t.WithFriendsFemaleFirstRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].FemaleRecommendations[i].IsWithFriends == constants.False {
			t.WithFriendsFemaleFirstRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}

	for i := 4; i < 8; i++ {
		if recommendations[gameId].FemaleRecommendations[i].IsWithFriends == constants.True {
			t.WithFriendsFemaleSecondRow |= uint8(int(math.Pow(2, float64(i))) << i)
		} else if recommendations[gameId].FemaleRecommendations[i].IsWithFriends == constants.False {
			t.WithFriendsFemaleSecondRow |= uint8(int(math.Pow(2, float64(i+1))) << i)
		}
	}
}

func (t *TitleTable) DetermineOtherFlags(game gametdb.Game) {
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
