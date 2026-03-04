package wc24data

import (
	"NintendoChannel/common"
	"NintendoChannel/constants"
	"NintendoChannel/gametdb"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"sync"

	"github.com/SketchMaster2001/libwc24crypt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wii-tools/lzx/lz10"
)

type WC24Data struct {
	Header        Header
	PlatformTable []Platform
	CompanyTable  []Company
	TitleTable    []Title
	NewTitleTable []uint32
	VideoTable    []VideoTable
	ExtraBytes    []uint16

	// State variables
	region      constants.Region
	ratingGroup constants.RatingGroup
	language    constants.Language
}

var (
	pool *pgxpool.Pool
	key  = []byte{17, 50, 20, 213, 122, 3, 143, 220, 230, 218, 224, 213, 173, 246, 135, 255}
	iv   = []byte{70, 70, 20, 40, 143, 110, 36, 6, 184, 107, 135, 239, 96, 45, 80, 151}
	ctx  = context.Background()
	// map[game_id]amount_voted
	recommendations = map[string]TitleRecommendation{}
)

func MakeList() {
	config := common.GetConfig()

	// Initialize database
	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", config.Username, config.Password, config.DatabaseAddress, config.DatabaseName)
	dbConf, err := pgxpool.ParseConfig(dbString)
	common.CheckError(err)
	pool, err = pgxpool.NewWithConfig(ctx, dbConf)
	common.CheckError(err)

	defer pool.Close()
	gametdb.PrepareGameTDB(config)
	PopulateRecommendations()

	wg := sync.WaitGroup{}
	semaphore := make(chan any, 3)

	wg.Add(10)
	for _, region := range constants.Regions {
		for _, language := range region.Languages {
			go func(_region constants.RegionMeta, _language constants.Language) {
				defer wg.Done()

				semaphore <- struct{}{}
				fmt.Printf("Starting worker - Region: %d, Language: %d\n", _region.Region, _language)
				list := WC24Data{
					region:      _region.Region,
					ratingGroup: _region.RatingGroup,
					language:    _language,
				}

				list.MakeHeader()
				list.MakePlatformTable()
				list.MakeCompaniesTable()
				list.MakeTitleTable()
				list.MakeNewTitleTable()
				list.MakeVideoTable()

				list.Header.Something = list.GetCurrentSize()
				// I am seriously confused. I think it has something to do with the data usage screen?
				list.ExtraBytes = make([]uint16, 6004/2)
				list.Header.Filesize = list.GetCurrentSize()

				buffer := new(bytes.Buffer)
				list.WriteAll(buffer)

				crcTable := crc32.MakeTable(crc32.IEEE)
				checksum := crc32.Checksum(buffer.Bytes(), crcTable)
				binary.BigEndian.PutUint32(buffer.Bytes()[8:], checksum)

				compress, err := lz10.Compress(buffer.Bytes())
				common.CheckError(err)

				rsaKey, err := os.ReadFile(fmt.Sprintf("%s/nc.pem", config.AssetsPath))
				common.CheckError(err)

				enc, err := libwc24crypt.EncryptWC24(compress, key, iv, rsaKey)
				common.CheckError(err)

				err = os.WriteFile(fmt.Sprintf("%s/lists/%d/%d/wc24data.LZ", config.AssetsPath, _region.Region, _language), enc, 0666)
				common.CheckError(err)

				fmt.Printf("Finished worker - Region: %d, Language: %d\n", _region.Region, _language)
				<-semaphore
			}(region, language)
		}
	}

	wg.Wait()
}

// Write writes the current values in Votes to an io.Writer method.
// This is required as Go cannot write structs with non-fixed slice sizes,
// but can write them individually.
func (w *WC24Data) Write(writer io.Writer, data any) {
	err := binary.Write(writer, binary.BigEndian, data)
	common.CheckError(err)
}

func (w *WC24Data) WriteAll(writer io.Writer) {
	w.Write(writer, w.Header)
	w.Write(writer, w.PlatformTable)
	w.Write(writer, w.CompanyTable)
	w.Write(writer, w.TitleTable)
	w.Write(writer, w.NewTitleTable)
	w.Write(writer, w.VideoTable)
	w.Write(writer, w.ExtraBytes)
}

// GetCurrentSize returns the current size of our List struct.
// This is useful for calculating the current offset of List.
func (w *WC24Data) GetCurrentSize() uint32 {
	buffer := bytes.NewBuffer(nil)
	w.WriteAll(buffer)

	return uint32(buffer.Len())
}
