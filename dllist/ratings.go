package dllist

import (
	"NintendoChannel/constants"
	"bytes"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	"unicode/utf16"
)

// RatingTable contains the data for the game ratings.
type RatingTable struct {
	RatingID    uint8
	RatingGroup constants.RatingGroup
	Age         uint8
	Unknown     uint8
	JPEGOffset  uint32
	JPEGSize    uint32
	RatingTitle [11]uint16
}

type DetailedRatingTable struct {
	RatingGroup constants.RatingGroup
	RatingID    uint8
	Title       [102]uint16
}

// MakeRatingsTable writes the rating levels for the current region.
func (l *List) MakeRatingsTable() {
	l.Header.RatingTableOffset = l.GetCurrentSize()

	for i, rating := range constants.RatingsData[l.ratingGroup] {
		ratingTable := RatingTable{
			RatingID:    uint8(i + 8),
			RatingGroup: l.ratingGroup,
			Age:         rating.Age,
			Unknown:     222,
			JPEGOffset:  0,
			JPEGSize:    0,
			RatingTitle: rating.Name,
		}

		l.RatingsTable = append(l.RatingsTable, ratingTable)
	}

	l.Header.NumberOfRatingTables = uint32(len(l.RatingsTable))
}

func (l *List) MakeDetailedRatingTable() {
	l.Header.DetailedRatingTablesOffset = l.GetCurrentSize()

	// TODO: Move away from kaitai
	dl := NewNinchDllist()
	err := dl.Read(kaitai.NewStream(bytes.NewReader(constants.DLList)), nil, dl)
	checkError(err)

	detailedRatings, err := dl.DetailedRatingsTable()
	checkError(err)

	for _, table := range detailedRatings {
		var title [102]uint16
		tempTitle := utf16.Encode([]rune(table.Title))
		copy(title[:], tempTitle)

		l.DetailedRatingTable = append(l.DetailedRatingTable, DetailedRatingTable{
			RatingGroup: constants.RatingGroup(table.RatingGroup),
			RatingID:    table.RatingId,
			Title:       title,
		})
	}

	l.Header.NumberOfDetailedRatingTables = uint32(len(l.DetailedRatingTable))
}

func (l *List) WriteRatingImages() {
	deadBeef := []byte{0xDE, 0xAD, 0xBE, 0xEF}

	for i, _ := range l.RatingsTable {
		l.RatingsTable[i].JPEGOffset = l.GetCurrentSize()
		l.RatingsTable[i].JPEGSize = uint32(len(constants.Images[l.ratingGroup][i]))
		l.imageBuffer.Write(constants.Images[l.ratingGroup][i])

		counter := 0
		for l.GetCurrentSize()%32 != 0 {
			l.imageBuffer.WriteByte(deadBeef[counter%4])
			counter++
		}
	}
}
