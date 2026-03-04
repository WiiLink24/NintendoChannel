package wc24data

type Header struct {
	_                       uint16
	Version                 uint8
	Unknown                 uint8
	Filesize                uint32
	CRC32                   uint32
	ListID                  uint32
	CountryCode             uint32
	LanguageCode            uint32
	SupportedLanguages      [16]byte
	Unknown1                [15]byte
	NumberOfPlatforms       uint32
	PlatformTableOffset     uint32
	NumberOfManufacturers   uint32
	ManufacturerTableOffset uint32
	NumberOfTitles          uint32
	TitleTableOffset        uint32
	NumberOfNewTitles       uint32
	NewTitleTableOffset     uint32
	NumberOfVideos          uint32
	VideoTableOffset        uint32
	NumberOfDemos           uint32
	DemoTableOffset         uint32
	Something               uint32
	DLListURLIDs            [768]byte
}

func (w *WC24Data) MakeHeader() {
	var tempUrlIDs []byte
	var urlIDs [768]byte
	for i := 0; i < 5; i++ {
		temp := make([]byte, 256)
		copy(temp[:64], "THqOxqSaiDd5bjhSQS6hk6nkYJVdioanD5Lc8mOHkobUkblWf8KxczDUZwY84FIV")
		tempUrlIDs = append(tempUrlIDs, temp...)
	}

	copy(urlIDs[:], tempUrlIDs)

	w.Header = Header{
		Version:     3,
		Unknown:     0,
		Filesize:    0,
		CRC32:       0,
		ListID:      981120054,
		CountryCode: 1,
		// This should be hardcoded or all the offsets somehow explode
		LanguageCode:            0,
		SupportedLanguages:      [16]byte{byte(w.language), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		Unknown1:                [15]byte{0, 1, 1, 0x51, 0xbe, 0xda, 0xaa, 0},
		NumberOfPlatforms:       0,
		PlatformTableOffset:     0,
		NumberOfManufacturers:   0,
		ManufacturerTableOffset: 0,
		NumberOfTitles:          0,
		TitleTableOffset:        0,
		NumberOfNewTitles:       0,
		NewTitleTableOffset:     0,
		NumberOfVideos:          0,
		VideoTableOffset:        0,
		NumberOfDemos:           0,
		DemoTableOffset:         0,
		Something:               0,
		DLListURLIDs:            urlIDs,
	}
}
