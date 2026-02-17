package constants

// RatingGroup is the rating organization for a specific region.
type RatingGroup uint8

const (
	_ RatingGroup = iota
	// CERO is the RatingGroup for NTSC-J games.
	CERO

	// ESRB is the RatingGroup for NTSC-U games.
	ESRB

	// PEGI is the RatingGroup for PAL games.
	PEGI RatingGroup = 4
)

// RatingData contains the name and age for a rating
type RatingData struct {
	Name [11]uint16
	Age  uint8
}

var RatingsData = map[RatingGroup][]RatingData{
	CERO: {
		{Name: [11]uint16{'A'}, Age: 0},
		{Name: [11]uint16{'B'}, Age: 12},
		{Name: [11]uint16{'C'}, Age: 15},
		{Name: [11]uint16{'D'}, Age: 17},
		{Name: [11]uint16{'Z'}, Age: 18},
	},
	ESRB: {
		{Name: [11]uint16{'E', 'C'}, Age: 3},
		{Name: [11]uint16{'E'}, Age: 6},
		{Name: [11]uint16{'E', '1', '0'}, Age: 10},
		{Name: [11]uint16{'T'}, Age: 13},
		{Name: [11]uint16{'M'}, Age: 17},
	},
	PEGI: {
		{Name: [11]uint16{'3'}, Age: 3},
		{Name: [11]uint16{'7'}, Age: 7},
		{Name: [11]uint16{'1', '2'}, Age: 12},
		{Name: [11]uint16{'1', '6'}, Age: 16},
		{Name: [11]uint16{'1', '8'}, Age: 18},
	},
}

// Region is the Wii's region flags found in TMDs.
type Region int

const (
	Japan Region = iota
	PAL
	NTSC
)

type Language int

const (
	Japanese = iota
	English
	German
	French
	Spanish
	Italian
	Dutch
)

type RegionMeta struct {
	Region      Region
	Languages   []Language
	RatingGroup RatingGroup
}

var Regions = []RegionMeta{
	{
		Region:      Japan,
		Languages:   []Language{Japanese},
		RatingGroup: CERO,
	},
	{
		Region:      NTSC,
		Languages:   []Language{English, French, Spanish},
		RatingGroup: ESRB,
	},
	{
		Region:      PAL,
		Languages:   []Language{English, German, French, Spanish, Italian, Dutch},
		RatingGroup: PEGI,
	},
}

// ConsoleModels is the type of consoles the Nintendo Channel games has.
type ConsoleModels [3]byte

var (
	// RVL represents titles on the Wii.
	RVL ConsoleModels = [3]byte{'R', 'V', 'L'}

	// NTR represents titles on the DS and DS Lite.
	NTR ConsoleModels = [3]byte{'N', 'T', 'R'}

	// TWL represents titles on the DSi.
	TWL ConsoleModels = [3]byte{'T', 'W', 'L'}

	// CTR represents titles on the 3DS.
	CTR ConsoleModels = [3]byte{'C', 'T', 'R'}
)

// TitleGroupTypes represents a type of title a console's game is.
type TitleGroupTypes uint8

const (
	_ TitleGroupTypes = iota

	// Disc represents Wii disc games.
	Disc

	// WiiWare represents WiiWare games.
	WiiWare

	// WiiChannels represents Wii Channels such as Forecast and News.
	WiiChannels

	// DS represents DS games.
	DS

	// VirtualConsole represents Virtual Console games.
	VirtualConsole

	// DSi represents games that support DSi only features
	DSi

	// DSiWare represents DSiWare games
	DSiWare

	// ThreeDS represents 3DS games
	ThreeDS

	// ThreeDSDownloadSoftware represents 3DS Download Software
	ThreeDSDownloadSoftware

	// ThreeDSGameBoy represents GameBoy Virtual Console games on the 3DS.
	ThreeDSGameBoy
)

// TitleTypeData contains the metadata for a title type
type TitleTypeData struct {
	ConsoleModel ConsoleModels
	GroupID      TitleGroupTypes
	TypeID       uint8
	ConsoleName  string
}

var TitleTypesData = []TitleTypeData{
	{TypeID: 1, ConsoleModel: RVL, ConsoleName: "Wii", GroupID: Disc},
	{TypeID: 11, ConsoleModel: RVL, ConsoleName: "WiiWare", GroupID: WiiWare},
	{TypeID: 2, ConsoleModel: RVL, ConsoleName: "Wii Channels", GroupID: WiiChannels},
	{TypeID: 3, ConsoleModel: RVL, ConsoleName: "NES", GroupID: VirtualConsole},
	{TypeID: 4, ConsoleModel: RVL, ConsoleName: "Super NES", GroupID: VirtualConsole},
	{TypeID: 5, ConsoleModel: RVL, ConsoleName: "Nintendo 64", GroupID: VirtualConsole},
	{TypeID: 6, ConsoleModel: RVL, ConsoleName: "TurboGrafx16", GroupID: VirtualConsole},
	{TypeID: 7, ConsoleModel: RVL, ConsoleName: "Sega Genesis", GroupID: VirtualConsole},
	{TypeID: 8, ConsoleModel: RVL, ConsoleName: "NEOGEO", GroupID: VirtualConsole},
	{TypeID: 12, ConsoleModel: RVL, ConsoleName: "MASTER SYSTEM", GroupID: VirtualConsole},
	{TypeID: 13, ConsoleModel: RVL, ConsoleName: "Commodore 64", GroupID: VirtualConsole},
	{TypeID: 14, ConsoleModel: RVL, ConsoleName: "Virtual Console Arcade", GroupID: VirtualConsole},
	{TypeID: 10, ConsoleModel: NTR, ConsoleName: "Nintendo DS", GroupID: DS},
	{TypeID: 15, ConsoleModel: TWL, ConsoleName: "Nintendo DS", GroupID: DS},
	{TypeID: 16, ConsoleModel: TWL, ConsoleName: "Nintendo DSi", GroupID: DSi},
	{TypeID: 17, ConsoleModel: TWL, ConsoleName: "Nintendo DSiWare", GroupID: DSiWare},
	{TypeID: 18, ConsoleModel: CTR, ConsoleName: "Nintendo 3DS", GroupID: ThreeDS},
	{TypeID: 19, ConsoleModel: CTR, ConsoleName: "Nintendo 3DS Download Software", GroupID: ThreeDSDownloadSoftware},
	{TypeID: 20, ConsoleModel: CTR, ConsoleName: "GAME BOY", GroupID: ThreeDSGameBoy},
}

var TouchGenIDs = []string{
	"YBN", "VAA", "AYA", "AND", "ANM", "ATD", "CVN",
	"YCU", "ATI", "AOS", "AG3", "AWI", "APL", "AJQ", "CM7",
	"AD5", "AD2", "ADG", "AD7", "AD3", "IMW", "C6P", "AXP",
	"A8N", "AZI", "ASQ", "ATR", "AGF",
	"RFN", "RFP", "R64", "RYW",
}

var PaynPlayIDs = []string{
	"WFC", "R3B", "WR9", "WRX", "SJD", "SD2", "SJX",
	"SJO", "SE3", "SZA", "SZB", "R9J", "SXE", "SXI", "R36",
	"SXA", "SWA", "SWB", "SXF", "R9O", "SUS", "SU3", "R83",
  
}
var DevAppIDs = []string{
	"007E", "091E", "410E", "413E", "5NEA", "RAAE",
}

// TitleType is the classified type of title according to GameTDB
type TitleType uint8

const (
	_ TitleType = iota
	Wii
	WiiChannel
	NES
	SNES
	Nintendo64
	TurboGrafx16
	Genesis
	NeoGeo
	NintendoDS           TitleType = 10
	WiiWare_             TitleType = 11
	MasterSystem         TitleType = 12
	Commodore64          TitleType = 13
	VirtualConsoleArcade TitleType = 14
	NintendoDSi          TitleType = 16
	NintendoDSiWare      TitleType = 17
	NintendoThreeDS      TitleType = 18
	ThreeDSDownload      TitleType = 19
)

var TitleTypeMap = map[string]TitleType{
	"Wii":       Wii,
	"Channel":   WiiChannel,
	"WiiWare":   WiiWare_,
	"VC-NES":    NES,
	"VC-SNES":   SNES,
	"VC-N64":    Nintendo64,
	"VC-SMS":    MasterSystem,
	"VC-MD":     Genesis,
	"VC-PCE":    TurboGrafx16,
	"VC-NEOGEO": NeoGeo,
	"VC-Arcade": VirtualConsoleArcade,
	"VC-C64":    Commodore64,
	"DS":        NintendoDS,
	"DSi":       NintendoDSi,
	"DSiWare":   NintendoDSiWare,
	"3DS":       NintendoThreeDS,
	"3DSWare":   ThreeDSDownload,
	"VC-GB":     ThreeDSDownload,
	"VC-GBC":    ThreeDSDownload,
	"VC-GBA":    ThreeDSDownload,
	"VC-GG":     ThreeDSDownload,
}

type Medal uint8

const (
	None Medal = iota
	Bronze
	Silver
	Gold
	Platinum
)

func GetVideoQueryString(language Language) string {
	switch language {
	case Japanese:
		return `SELECT id, name_japanese, length, video_type, date_added FROM videos ORDER BY id DESC`
	case English:
		return `SELECT id, name_english, length, video_type, date_added FROM videos ORDER BY id DESC`
	case German:
		return `SELECT id, name_german, length, video_type, date_added FROM videos ORDER BY id DESC`
	case French:
		return `SELECT id, name_french, length, video_type, date_added FROM videos ORDER BY id DESC`
	case Spanish:
		return `SELECT id, name_spanish, length, video_type, date_added FROM videos ORDER BY id DESC`
	case Italian:
		return `SELECT id, name_italian, length, video_type, date_added FROM videos ORDER BY id DESC`
	case Dutch:
		return `SELECT id, name_dutch, length, video_type, date_added FROM videos ORDER BY id DESC`
	default:
		// Will never reach here
		return ""
	}
}

// CriteriaBool is a boolean like value for AgeRecommendationData criteria.
// A criteria may either be true, false, or none depending on if the titles has any recommendations.
type CriteriaBool int

const (
	False CriteriaBool = iota
	True
	Nil
)

type AgeRecommendationData struct {
	LowerAge        int
	UpperAge        int
	IsGamers        CriteriaBool
	IsHardcore      CriteriaBool
	IsWithFriends   CriteriaBool
	EveryonePercent uint8
	CasualPercent   uint8
	AlonePercent    uint8
}

// AgeRecommendationTable exists with the sole purpose of populating the age bounds.
var AgeRecommendationTable = [8]AgeRecommendationData{
	{0, 100, Nil, Nil, Nil, 0, 0, 0},
	{0, 12, Nil, Nil, Nil, 0, 0, 0},
	{13, 18, Nil, Nil, Nil, 0, 0, 0},
	{19, 24, Nil, Nil, Nil, 0, 0, 0},
	{25, 34, Nil, Nil, Nil, 0, 0, 0},
	{35, 44, Nil, Nil, Nil, 0, 0, 0},
	{45, 54, Nil, Nil, Nil, 0, 0, 0},
	{55, 100, Nil, Nil, Nil, 0, 0, 0},
}
