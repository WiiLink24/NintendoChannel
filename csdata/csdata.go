package csdata

import (
	"NintendoChannel/common"
	"NintendoChannel/constants"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/SketchMaster2001/libwc24crypt"
	"github.com/jackc/pgx/v4/pgxpool"
	tpl "github.com/wii-tools/libtpl"
	"github.com/wii-tools/lzx/lz10"
	"hash/crc32"
	"image/jpeg"
	"log"
	"os"
	"unicode/utf16"
)

type DBBanner struct {
	ID       int
	Japanese string
	English  string
	German   string
	French   string
	Spanish  string
	Italian  string
	Dutch    string
}

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

func (d *DBBanner) GetTextForLanguage(language constants.Language) string {
	switch language {
	case constants.Japanese:
		return d.Japanese
	case constants.English:
		return d.English
	case constants.German:
		return d.German
	case constants.French:
		return d.French
	case constants.Spanish:
		return d.Spanish
	case constants.Italian:
		return d.Italian
	case constants.Dutch:
		return d.Dutch
	default:
		// Will never reach here
		return ""
	}
}

func CreateCSData() {
	ctx := context.Background()
	config := common.GetConfig()

	// Initialize database
	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", config.Username, config.Password, config.DatabaseAddress, config.DatabaseName)
	dbConf, err := pgxpool.ParseConfig(dbString)
	common.CheckError(err)
	pool, err := pgxpool.ConnectConfig(ctx, dbConf)
	common.CheckError(err)

	// Ensure this Postgresql connection is valid.
	defer pool.Close()

	rows, err := pool.Query(ctx, `SELECT * FROM banners`)
	common.CheckError(err)

	var dbBanners []DBBanner
	for rows.Next() {
		var dbBanner DBBanner
		err = rows.Scan(&dbBanner.ID, &dbBanner.Japanese, &dbBanner.English, &dbBanner.German, &dbBanner.French,
			&dbBanner.Spanish, &dbBanner.Italian, &dbBanner.Dutch)
		common.CheckError(err)

		dbBanners = append(dbBanners, dbBanner)
	}

	// There should be 3 banners.
	// TODO: Figure out if more or less can be added.
	if len(dbBanners) != 3 {
		log.Fatalf("Expected 3 rows in banners! Got %d\n", len(dbBanners))
	}

	var pics [][]byte
	for _, banner := range dbBanners {
		data, err := os.Open(fmt.Sprintf("%s/banners/%d.img", config.AssetsPath, banner.ID))
		common.CheckError(err)

		pic, err := jpeg.Decode(data)
		common.CheckError(err)

		// Saved as a JPEG image. Convert to TPL for the Wii to read.
		enc, err := tpl.ToRGB565(pic)
		common.CheckError(err)

		// Append to slice while stripping the TPL header.
		pics = append(pics, enc[64:])
	}

	// First append the DLListID to a
	var DLListID [256]byte
	tempID := make([]byte, 256)
	copy(tempID[:65], "6THqOxqSaiDd5bjhSQS6hk6nkYJVdioanD5Lc8mOHkobUkblWf8KxczDUZwY84FIV")
	copy(DLListID[:], tempID)

	// Generate for all regions and languages
	for _, region := range constants.Regions {
		for _, language := range region.Languages {
			languagesRaw := [16]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
			// Can't cast the entire languages slice to u8 so for loop it is
			for i, l := range region.Languages {
				languagesRaw[i] = byte(l)
			}

			header := Header{
				Version:            6,
				Unknown:            0,
				Filesize:           0,
				CRC32:              0,
				DLListID:           434968891,
				CountryCode:        49,
				LanguageCode:       uint32(language),
				SupportedLanguages: languagesRaw,
				Unknown1:           [12]byte{0, 78, 112, 38, 194, 0, 0, 0, 3, 0, 0, 1},
				DLUrlID:            DLListID,
				Unknown2:           222,
			}

			for i := 0; i < 3; i++ {
				var textArray [51]uint16
				tempText := utf16.Encode([]rune(dbBanners[i].GetTextForLanguage(language)))
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

			err = binary.Write(buffer, binary.BigEndian, header)
			err = binary.Write(buffer, binary.BigEndian, pics[0])
			err = binary.Write(buffer, binary.BigEndian, pics[1])
			err = binary.Write(buffer, binary.BigEndian, pics[2])
			common.CheckError(err)

			header.Filesize = uint32(buffer.Len())
			buffer.Reset()

			err = binary.Write(buffer, binary.BigEndian, header)
			err = binary.Write(buffer, binary.BigEndian, pics[0])
			err = binary.Write(buffer, binary.BigEndian, pics[1])
			err = binary.Write(buffer, binary.BigEndian, pics[2])
			common.CheckError(err)

			// Calculate crc32
			crcTable := crc32.MakeTable(crc32.IEEE)
			checksum := crc32.Checksum(buffer.Bytes(), crcTable)
			header.CRC32 = checksum
			buffer.Reset()

			err = binary.Write(buffer, binary.BigEndian, header)
			err = binary.Write(buffer, binary.BigEndian, pics[0])
			err = binary.Write(buffer, binary.BigEndian, pics[1])
			err = binary.Write(buffer, binary.BigEndian, pics[2])
			common.CheckError(err)

			compress, err := lz10.Compress(buffer.Bytes())
			common.CheckError(err)

			rsaKey, err := os.ReadFile(fmt.Sprintf("%s/nc.pem", config.AssetsPath))
			common.CheckError(err)

			encrypted, err := libwc24crypt.EncryptWC24(compress, key, iv, rsaKey)
			common.CheckError(err)

			// Create directory just in case.
			err = os.MkdirAll(fmt.Sprintf("%s/csdata/%d/%d", config.AssetsPath, region.Region, language), 0777)
			common.CheckError(err)

			err = os.WriteFile(fmt.Sprintf("%s/csdata/%d/%d/csdata.bin", config.AssetsPath, region.Region, language), encrypted, 0777)
			common.CheckError(err)
		}
	}
}
