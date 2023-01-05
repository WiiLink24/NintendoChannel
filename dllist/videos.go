package dllist

import (
	"unicode/utf16"
)

type VideoTable struct {
	ID          uint32
	VideoLength uint16
	TitleID     uint32
	VideoType   uint8
	Unknown     [14]byte
	Unknown2    uint8
	RatingID    uint8
	Unknown3    uint8
	IsNew       uint8
	// Starts at 1 and incremented by 1
	VideoIndex uint8
	Unknown4   [2]byte
	Title      [123]uint16
}

type NewVideoTable struct {
	ID          uint32
	VideoLength uint16
	TitleID     uint32
	Unknown     [15]byte
	Unknown2    uint8
	RatingID    uint8
	Unknown3    uint8
	Title       [102]uint16
}

type PopularVideosTable struct {
	ID          uint32
	VideoLength uint16
	TitleID     uint32
	BarColor    uint8
	_           [15]byte
	RatingID    uint8
	Unknown     uint8
	VideoRank   uint8
	Unknown2    uint8
	Title       [102]uint16
}

func (l *List) MakeVideoTable() {
	l.Header.VideoTableOffset = l.GetCurrentSize()

	rows, err := pool.Query(ctx, `SELECT name_english, length, video_type FROM videos ORDER BY id ASC`)
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		var queriedTitle string
		var length int
		var videoType int

		err = rows.Scan(&queriedTitle, &length, &videoType)
		checkError(err)

		var title [123]uint16
		tempTitle := utf16.Encode([]rune(queriedTitle))
		copy(title[:], tempTitle)

		l.VideoTable = append(l.VideoTable, VideoTable{
			ID:          1,
			VideoLength: uint16(length),
			TitleID:     0,
			VideoType:   uint8(videoType),
			Unknown:     [14]byte{},
			Unknown2:    0,
			RatingID:    9,
			Unknown3:    1,
			IsNew:       1,
			VideoIndex:  1,
			Unknown4:    [2]byte{61, 222},
			Title:       title,
		})
	}

	l.Header.NumberOfVideoTables = uint32(len(l.VideoTable))
}

func (l *List) MakeNewVideoTable() {
	l.Header.NewVideoTableOffset = l.GetCurrentSize()

	var title [102]uint16
	tempTitle := utf16.Encode([]rune("WiiLink goes global!"))
	copy(title[:], tempTitle)

	l.NewVideoTable = append(l.NewVideoTable, NewVideoTable{
		ID:          1,
		VideoLength: 280,
		TitleID:     0,
		Unknown:     [15]byte{},
		Unknown2:    0,
		RatingID:    9,
		Unknown3:    1,
		Title:       title,
	})
	l.Header.NumberOfNewVideoTables = 1
}

func (l *List) MakePopularVideoTable() {
	l.Header.PopularVideoTableOffset = l.GetCurrentSize()

	var title [102]uint16
	tempTitle := utf16.Encode([]rune("WiiLink goes global!"))
	copy(title[:], tempTitle)

	l.PopularVideosTable = append(l.PopularVideosTable, PopularVideosTable{
		ID:          1,
		VideoLength: 280,
		TitleID:     0,
		BarColor:    0,
		RatingID:    9,
		Unknown:     1,
		VideoRank:   1,
		Unknown2:    222,
		Title:       title,
	})
	l.Header.NumberOfPopularVideoTables = 1
}
