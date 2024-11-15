package dllist

import (
	"NintendoChannel/common"
	"NintendoChannel/constants"
	"NintendoChannel/gametdb"
	"NintendoChannel/info"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/wii-tools/lzx/lz10"
	"hash/crc32"
	"io"
	"os"
	"runtime"
	"sync"
)

type List struct {
	Header          Header
	RatingsTable    []RatingTable
	TitleTypesTable []TitleTypeTable
	CompaniesTable  []CompanyTable
	TitleTable      []TitleTable
	// NewTitleTable is an array of pointers to titles in TitleTable
	NewTitleTable             []uint32
	VideoTable                []VideoTable
	NewVideoTable             []NewVideoTable
	DemoTable                 []DemoTable
	RecommendationTable       []uint32
	RecentRecommendationTable []RecentRecommendationTable
	PopularVideosTable        []PopularVideosTable
	DetailedRatingTable       []DetailedRatingTable

	// Below are variables that help us keep state
	region      constants.Region
	ratingGroup constants.RatingGroup
	language    constants.Language
	// map[game_id]amount_voted
	recommendations map[string]TitleRecommendation
	imageBuffer     *bytes.Buffer
}

var (
	config         common.Config
	pool           *pgxpool.Pool
	ctx            = context.Background()
	generateTitles = true
)

func MakeDownloadList(_generateTitles bool) {
	generateTitles = _generateTitles

	config = common.GetConfig()

	// Initialize database
	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", config.Username, config.Password, config.DatabaseAddress, config.DatabaseName)
	dbConf, err := pgxpool.ParseConfig(dbString)
	common.CheckError(err)
	pool, err = pgxpool.ConnectConfig(ctx, dbConf)
	common.CheckError(err)

	// Ensure this Postgresql connection is valid.
	defer pool.Close()
	gametdb.PrepareGameTDB()
	info.GetTimePlayed(&ctx, pool)

	wg := sync.WaitGroup{}
	runtime.GOMAXPROCS(runtime.NumCPU())
	semaphore := make(chan struct{}, 3)

	wg.Add(10)
	for _, region := range constants.Regions {
		for _, language := range region.Languages {
			go func(_region constants.RegionMeta, _language constants.Language) {
				defer wg.Done()
				if _region.Region != constants.NTSC || _language != constants.English {
					return
				}

				semaphore <- struct{}{}
				fmt.Printf("Starting worker - Region: %d, Language: %d\n", _region.Region, _language)
				list := List{
					region:          _region.Region,
					ratingGroup:     _region.RatingGroup,
					language:        _language,
					imageBuffer:     new(bytes.Buffer),
					recommendations: map[string]TitleRecommendation{},
				}

				list.QueryRecommendations()

				list.MakeHeader()
				list.MakeRatingsTable()
				list.MakeTitleTypeTable()
				list.MakeCompaniesTable()
				list.MakeTitleTable()
				list.MakeNewTitleTable()
				list.MakeVideoTable()
				list.MakeNewVideoTable()
				list.MakeDemoTable()
				list.MakeRecommendationTable()
				list.MakeRecentRecommendationTable()
				list.MakePopularVideoTable()
				list.MakeDetailedRatingTable()
				list.WriteRatingImages()

				temp := bytes.NewBuffer(nil)
				list.WriteAll(temp)
				list.Header.Filesize = uint32(temp.Len())
				temp.Reset()
				list.WriteAll(temp)

				crcTable := crc32.MakeTable(crc32.IEEE)
				checksum := crc32.Checksum(temp.Bytes(), crcTable)
				list.Header.CRC32 = checksum

				temp.Reset()
				list.WriteAll(temp)

				// Compress then write
				compressed, err := lz10.Compress(temp.Bytes())
				common.CheckError(err)

				err = os.WriteFile(fmt.Sprintf("%s/lists/%d/%d/dllist.bin", config.AssetsPath, _region.Region, _language), compressed, 0666)
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
func (l *List) Write(writer io.Writer, data any) {
	err := binary.Write(writer, binary.BigEndian, data)
	common.CheckError(err)
}

func (l *List) WriteAll(writer io.Writer) {
	l.Write(writer, l.Header)
	l.Write(writer, l.RatingsTable)
	l.Write(writer, l.TitleTypesTable)
	l.Write(writer, l.CompaniesTable)
	l.Write(writer, l.TitleTable)
	l.Write(writer, l.NewTitleTable)
	l.Write(writer, l.VideoTable)
	l.Write(writer, l.NewVideoTable)
	l.Write(writer, l.DemoTable)
	l.Write(writer, l.RecommendationTable)
	l.Write(writer, l.RecentRecommendationTable)
	l.Write(writer, l.PopularVideosTable)
	l.Write(writer, l.DetailedRatingTable)
}

// GetCurrentSize returns the current size of our List struct.
// This is useful for calculating the current offset of List.
func (l *List) GetCurrentSize() uint32 {
	buffer := bytes.NewBuffer(nil)
	l.WriteAll(buffer)
	buffer.Write(l.imageBuffer.Bytes())

	return uint32(buffer.Len())
}
