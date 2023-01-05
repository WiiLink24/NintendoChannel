package dllist

import (
	"unicode/utf16"
)

// Header represents the structure of the Header of dllist.bin
type Header struct {
	_                                  uint16
	Version                            uint8
	Region                             uint8
	Filesize                           uint32
	CRC32                              uint32
	ListID                             uint32
	ThumbnailID                        uint32
	CountryCode                        uint32
	LanguageCode                       uint32
	UnknownValue                       [9]byte
	NumberOfRatingTables               uint32
	RatingTableOffset                  uint32
	NumberOfTitleTypeTables            uint32
	TitleTypeTableOffset               uint32
	NumberOfCompanyTables              uint32
	CompanyTableOffset                 uint32
	NumberOfTitleTables                uint32
	TitleTableOffset                   uint32
	NumberOfNewTitleTables             uint32
	NewTitleTableOffset                uint32
	NumberOfVideoTables                uint32
	VideoTableOffset                   uint32
	NumberOfNewVideoTables             uint32
	NewVideoTableOffset                uint32
	NumberOfDemoTables                 uint32
	DemoTableOffset                    uint32
	_                                  uint64
	NumberOfRecommendationTables       uint32
	RecommendationTableOffset          uint32
	_                                  [2]uint64
	NumberOfRecentRecommendationTables uint32
	RecentRecommendationTableOffset    uint32
	_                                  uint64
	NumberOfPopularVideoTables         uint32
	PopularVideoTableOffset            uint32
	NumberOfDetailedRatingTables       uint32
	DetailedRatingTablesOffset         uint32
	LastUpdate                         [31]uint16
	UnknownValue2                      [3]byte
	DlUrlIDs                           [1280]byte
	UnknownValue3                      uint32
}

func (l *List) MakeHeader() {
	// First format the update name to UTF-16.
	var update [31]uint16
	tempUpdate := utf16.Encode([]rune("WiiLink Edition"))
	copy(update[:], tempUpdate)

	// Next write the DlUrlIDs
	var tempUrlIDs []byte
	var urlIDs [1280]byte
	for i := 0; i < 5; i++ {
		temp := make([]byte, 256)
		copy(temp[:64], "THqOxqSaiDd5bjhSQS6hk6nkYJVdioanD5Lc8mOHkobUkblWf8KxczDUZwY84FIV")
		tempUrlIDs = append(tempUrlIDs, temp...)

	}

	copy(urlIDs[:], tempUrlIDs)

	l.Header = Header{
		Version:                            6,
		Region:                             2,
		Filesize:                           0,
		CRC32:                              0,
		ListID:                             340086107,
		ThumbnailID:                        1,
		CountryCode:                        18,
		LanguageCode:                       uint32(l.language),
		UnknownValue:                       [9]byte{1, 0x50, 0x3C, 0xEF, 0, 0, 0, 0, 0},
		NumberOfRatingTables:               0,
		RatingTableOffset:                  0,
		NumberOfTitleTypeTables:            0,
		TitleTypeTableOffset:               0,
		NumberOfCompanyTables:              0,
		CompanyTableOffset:                 0,
		NumberOfTitleTables:                0,
		TitleTableOffset:                   0,
		NumberOfNewTitleTables:             0,
		NewTitleTableOffset:                0,
		NumberOfVideoTables:                0,
		VideoTableOffset:                   0,
		NumberOfNewVideoTables:             0,
		NewVideoTableOffset:                0,
		NumberOfDemoTables:                 0,
		DemoTableOffset:                    0,
		NumberOfRecommendationTables:       0,
		RecommendationTableOffset:          0,
		NumberOfRecentRecommendationTables: 0,
		RecentRecommendationTableOffset:    0,
		NumberOfPopularVideoTables:         0,
		PopularVideoTableOffset:            0,
		NumberOfDetailedRatingTables:       0,
		DetailedRatingTablesOffset:         0,
		LastUpdate:                         update,
		UnknownValue2:                      [3]byte{0, 0x2E, 1},
		DlUrlIDs:                           urlIDs,
		UnknownValue3:                      285278430,
	}
}
