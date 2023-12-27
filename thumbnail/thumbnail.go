package thumbnail

import (
	"NintendoChannel/constants"
	"bytes"
	"context"
	_ "embed"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	colorFmt "github.com/fatih/color"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Thumbnail struct {
	_              uint16
	Version        uint8
	Unknown        uint8
	Filesize       uint32
	Unknown1       uint32
	LanguageCode   uint32
	CountryCode    uint32
	Unknown2       uint32
	Unknown3       uint32
	NumberOfImages uint32
}

type ImageTable struct {
	ImageSize   uint32
	ImageOffset uint32
}

const ThumbnailHeaderSize = 32

// Make text bold
func bold(text string) string {
	return "\033[1m" + text + "\033[0m"
}

func checkError(err error) {
	if err != nil {
		// ERROR! bold and red
		colorFmt.HiRed(bold("An error has occurred!"))
		fmt.Println()
		log.Fatalf(bold("Nintendo Channel file generator has encountered a fatal error!\n\n" + bold("Reason: ") + err.Error() + "\n"))
	}
}

// Database credentials (you'll need to change these for your own database)
// Learn how to set up a PostgreSQL database here: https://www.postgresql.org/docs/13/tutorial-start.html
const (
	dbUser     = "user"
	dbPassword = "password"
	dbHost     = "127.0.0.1"
	dbName     = "nintendochannel"
)

func WriteThumbnail() {
	// Initialize database
	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", dbUser, dbPassword, dbHost, dbName)

	dbConf, err := pgxpool.ParseConfig(dbString)
	checkError(err)

	pool, err := pgxpool.ConnectConfig(context.Background(), dbConf)
	checkError(err)

	// Ensure this Postgresql connection is valid.
	defer pool.Close()

	rows, err := pool.Query(context.Background(), constants.GetVideoQueryString(constants.English))
	checkError(err)

	var images []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id, nil, nil, nil, nil)
		checkError(err)
		images = append(images, id)
	}

	buffer := new(bytes.Buffer)
	imageBuffer := new(bytes.Buffer)
	header := Thumbnail{
		Version:        6,
		Unknown:        2,
		Filesize:       0,
		Unknown1:       601820255,
		LanguageCode:   1,
		CountryCode:    49,
		Unknown2:       1,
		Unknown3:       1252951207,
		NumberOfImages: uint32(len(images) * 2),
	}

	err = binary.Write(buffer, binary.BigEndian, header)
	checkError(err)

	deadBeef := []byte{0xDE, 0xAD, 0xBE, 0xEF}

	for _, image := range images {
		file, err := os.ReadFile(fmt.Sprintf("/path/to/videos/%d.img", image))
		checkError(err)

		table := ImageTable{
			ImageSize:   uint32(len(file)),
			ImageOffset: uint32((ThumbnailHeaderSize + 8*(len(images)*2)) + imageBuffer.Len()),
		}

		err = binary.Write(buffer, binary.BigEndian, table)
		checkError(err)

		imageBuffer.Write(file)

		counter := 0
		for (256+imageBuffer.Len())%32 != 0 {
			imageBuffer.WriteByte(deadBeef[counter%4])
			counter++
		}
	}

	// Write twice because yes
	for _, image := range images {
		file, err := os.ReadFile(fmt.Sprintf("/path/to/videos/%d.img", image))
		checkError(err)

		table := ImageTable{
			ImageSize:   uint32(len(file)),
			ImageOffset: uint32((ThumbnailHeaderSize + 8*(len(images)*2)) + imageBuffer.Len()),
		}

		err = binary.Write(buffer, binary.BigEndian, table)
		checkError(err)

		imageBuffer.Write(file)

		counter := 0
		for (256+imageBuffer.Len())%32 != 0 {
			imageBuffer.WriteByte(deadBeef[counter%4])
			counter++
		}
	}

	buffer.Write(imageBuffer.Bytes())
	binary.BigEndian.PutUint32(buffer.Bytes()[4:8], uint32(buffer.Len()))

	err = os.WriteFile("thumbnail.bin", buffer.Bytes(), 0666)
	checkError(err)
}
