package csdata

import (
	"bytes"
	"encoding/binary"
	"github.com/SketchMaster2001/libwc24crypt"
	"github.com/wii-tools/lzx/lz10"
	"hash/crc32"
	"os"
	"unicode/utf16"
)

type Header struct {
	_                  uint16
	Version            uint8
	Unknown            uint8
	Filesize           uint32
	CRC32              uint32
	DLListID           uint32
	CountryCode        uint32
	LanguageCode       uint32
	SupportedLanguages [16]byte
	Unknown1           [12]byte
	DLUrlID            [256]byte
	Unknown2           uint16
	Banners            [3]Banner
}

type Banner struct {
	Text          [51]uint16
	PictureSize   uint32
	PictureOffset uint32
}

var (
	key = []byte{17, 50, 20, 213, 122, 3, 143, 220, 230, 218, 224, 213, 173, 246, 135, 255}
	iv  = []byte{70, 70, 20, 40, 143, 110, 36, 6, 184, 107, 135, 239, 96, 45, 80, 151}
)

func CreateCSData() {
	// First append the DLListID to a
	var DLListID [256]byte
	tempID := make([]byte, 256)
	copy(tempID[:65], "6THqOxqSaiDd5bjhSQS6hk6nkYJVdioanD5Lc8mOHkobUkblWf8KxczDUZwY84FIV")
	copy(DLListID[:], tempID)

	header := Header{
		Version:            6,
		Unknown:            2,
		Filesize:           0,
		CRC32:              0,
		DLListID:           434968891,
		CountryCode:        49,
		LanguageCode:       3,
		SupportedLanguages: [16]byte{1, 2, 3, 4, 5, 6, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		Unknown1:           [12]byte{0, 78, 112, 38, 194, 0, 0, 0, 3, 0, 0, 1},
		DLUrlID:            DLListID,
		Unknown2:           222,
	}

	var pics [][]byte
	pic1, err := os.ReadFile("pic1.tpl")
	if err != nil {
		panic(err)
	}

	pic2, err := os.ReadFile("pic2.tpl")
	if err != nil {
		panic(err)
	}

	pic3, err := os.ReadFile("pic1.tpl")
	if err != nil {
		panic(err)
	}

	pics = append(pics, pic1[64:], pic2[64:], pic3[64:])

	text := []string{"WiiLink Nintendo Channel!", "Demae Channel released", "WiiLink goes global"}
	for i := 0; i < 3; i++ {
		var textArray [51]uint16
		tempText := utf16.Encode([]rune(text[i]))
		copy(textArray[:], tempText)

		var offset uint32
		if i == 0 {
			offset = 640
		} else if i == 1 {
			offset = uint32(640 + len(pics[0]))
		} else if i == 2 {
			offset = uint32(640 + len(pics[0]) + len(pics[1]))
		}

		header.Banners[i] = Banner{
			Text:          textArray,
			PictureSize:   uint32(len(pics[i])),
			PictureOffset: offset,
		}
	}

	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, header)
	binary.Write(buffer, binary.BigEndian, pics[0])
	binary.Write(buffer, binary.BigEndian, pics[1])
	binary.Write(buffer, binary.BigEndian, pics[2])
	header.Filesize = uint32(buffer.Len())
	buffer.Reset()

	binary.Write(buffer, binary.BigEndian, header)
	binary.Write(buffer, binary.BigEndian, pics[0])
	binary.Write(buffer, binary.BigEndian, pics[1])
	binary.Write(buffer, binary.BigEndian, pics[2])

	// Calculate crc32
	crcTable := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum(buffer.Bytes(), crcTable)
	header.CRC32 = checksum
	buffer.Reset()

	binary.Write(buffer, binary.BigEndian, header)
	binary.Write(buffer, binary.BigEndian, pics[0])
	binary.Write(buffer, binary.BigEndian, pics[1])
	binary.Write(buffer, binary.BigEndian, pics[2])

	compress, err := lz10.Compress(buffer.Bytes())

	rsaKey, err := os.ReadFile("nc.pem")
	encrypted, err := libwc24crypt.EncryptWC24(compress, key, iv, rsaKey)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("csdata.bin", encrypted, 0666)
	if err != nil {
		panic(err)
	}
}
