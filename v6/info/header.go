package info

import "NintendoChannel/constants"

type Header struct {
	_                                   uint16
	Version                             uint8
	Unknown                             uint8
	Filesize                            uint32
	CRC32                               uint32
	DLListID                            uint32
	CountryCode                         uint32
	LanguageCode                        uint32
	RatingTableOffset                   uint32
	TimesPlayedTableOffset              uint32
	NumberOfPeopleWhoLikedThisAlsoLiked uint32
	PeopleWhoLikedThisAlsoLikedOffset   uint32
	NumberOfRelatedTitlesTables         uint32
	RelatedTitlesTableOffset            uint32
	NumberOfVideoTables                 uint32
	VideoTableOffset                    uint32
	NumberOfDemosTables                 uint32
	DemosTableOffset                    uint32
	_                                   uint64
	PictureOffset                       uint32
	PictureSize                         uint32
	_                                   uint32
	RatingPictureOffset                 uint32
	RatingPictureSize                   uint32
	DetailedRatingPictureTable          [7]DetailedRatingPictureTable
	_                                   uint64
	SoftwareID                          uint32
	GameID                              [4]byte
	TitleType                           constants.TitleType
	CompanyID                           uint32
	Unknown2                            uint16
	Unknown3                            uint16
	Unknown4                            uint8
	IsOnWiiShop                         uint8
	IsPurchasable                       uint8
	ReleaseYear                         uint16
	ReleaseMonth                        uint8
	ReleaseDay                          uint8
	ShopPoints                          uint32
	Unknown5                            [3]byte
	NumberOfPlayers                     uint8
}

func (i *Info) MakeHeader(softwareID uint32, gameID [4]byte, numberOfPlayers uint8, companyID uint32, titleType constants.TitleType, releaseYear uint16, releaseMonth, releaseDay uint8) {
	i.Header = Header{
		Version:                             6,
		Unknown:                             2,
		Filesize:                            0,
		CRC32:                               0,
		DLListID:                            1254762001,
		CountryCode:                         49,
		LanguageCode:                        1,
		RatingTableOffset:                   0,
		TimesPlayedTableOffset:              0,
		NumberOfPeopleWhoLikedThisAlsoLiked: 0,
		PeopleWhoLikedThisAlsoLikedOffset:   0,
		NumberOfRelatedTitlesTables:         0,
		RelatedTitlesTableOffset:            0,
		NumberOfVideoTables:                 0,
		VideoTableOffset:                    0,
		NumberOfDemosTables:                 0,
		DemosTableOffset:                    0,
		PictureOffset:                       0,
		PictureSize:                         0,
		RatingPictureOffset:                 0,
		RatingPictureSize:                   0,
		DetailedRatingPictureTable:          [7]DetailedRatingPictureTable{},
		SoftwareID:                          softwareID,
		GameID:                              gameID,
		TitleType:                           titleType,
		CompanyID:                           companyID,
		Unknown2:                            1,
		Unknown3:                            1,
		Unknown4:                            1,
		IsOnWiiShop:                         0,
		IsPurchasable:                       0,
		ReleaseYear:                         releaseYear,
		ReleaseMonth:                        releaseMonth,
		ReleaseDay:                          releaseDay,
		ShopPoints:                          0xFFFFFFFF,
		Unknown5:                            [3]byte{4, 1, 0},
		NumberOfPlayers:                     numberOfPlayers,
	}
	switch titleType {
	case constants.WiiWare_,
		constants.NES,
		constants.SNES,
		constants.Nintendo64,
		constants.TurboGrafx16,
		constants.Genesis,
		constants.NeoGeo,
		constants.MasterSystem,
		constants.Commodore64,
		constants.VirtualConsoleArcade:
		i.Header.IsOnWiiShop = 1
	default:
		i.Header.IsOnWiiShop = 0
	}
	switch titleType {
	case constants.Wii,
		constants.NintendoDS,
		constants.NintendoDSi,
		constants.NintendoThreeDS:
		i.Header.IsPurchasable = 1
	default:
		i.Header.IsPurchasable = 0
	}
	switch titleType { // Can't guarantee the amount of shop points for anything but Virtual Console
	case constants.NES:
		i.Header.ShopPoints = 500
	case constants.SNES:
		i.Header.ShopPoints = 800
	case constants.Nintendo64:
		i.Header.ShopPoints = 1000
	case constants.Genesis:
		i.Header.ShopPoints = 800
	case constants.NeoGeo:
		i.Header.ShopPoints = 900
	case constants.MasterSystem:
		i.Header.ShopPoints = 500
	default:
		i.Header.ShopPoints = 0xFFFFFFFF
	} //TODO: implement the rest of the shop points for different regions (imports may be 100 or more points than usa) and publishers cuz inconsistency
}
