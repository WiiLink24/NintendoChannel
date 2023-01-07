package dllist

import (
	"NintendoChannel/constants"
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

func (l *List) getVideoQueryString() string {
	switch l.language {
	case constants.Japanese:
		return `SELECT id, name_japanese, length, video_type FROM videos ORDER BY id DESC`
	case constants.English:
		return `SELECT id, name_english, length, video_type FROM videos ORDER BY id DESC`
	case constants.German:
		return `SELECT id, name_german, length, video_type FROM videos ORDER BY id DESC`
	case constants.French:
		return `SELECT id, name_french, length, video_type FROM videos ORDER BY id DESC`
	case constants.Spanish:
		return `SELECT id, name_spanish, length, video_type FROM videos ORDER BY id DESC`
	case constants.Italian:
		return `SELECT id, name_italian, length, video_type FROM videos ORDER BY id DESC`
	case constants.Dutch:
		return `SELECT id, name_japanese, length, video_type FROM videos ORDER BY id DESC`
	default:
		// Will never reach here
		return ""
	}
}

func (l *List) MakeVideoTable() {
	l.Header.VideoTableOffset = l.GetCurrentSize()

	rows, err := pool.Query(ctx, l.getVideoQueryString())
	checkError(err)

	index := 1
	defer rows.Close()
	for rows.Next() {
		var id int
		var queriedTitle string
		var length int
		var videoType int

		err = rows.Scan(&id, &queriedTitle, &length, &videoType)
		checkError(err)

		var title [123]uint16
		tempTitle := utf16.Encode([]rune(queriedTitle))
		copy(title[:], tempTitle)

		l.VideoTable = append(l.VideoTable, VideoTable{
			ID:          uint32(id),
			VideoLength: uint16(length),
			TitleID:     0,
			VideoType:   uint8(videoType),
			Unknown:     [14]byte{},
			Unknown2:    0,
			RatingID:    9,
			Unknown3:    1,
			IsNew:       1,
			VideoIndex:  uint8(index),
			Unknown4:    [2]byte{222, 222},
			Title:       title,
		})
		index++
	}

	l.Header.NumberOfVideoTables = uint32(len(l.VideoTable))
}

func (l *List) MakeNewVideoTable() {
	l.Header.NewVideoTableOffset = l.GetCurrentSize()

	rows, err := pool.Query(ctx, l.getVideoQueryString())
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		var id int
		var queriedTitle string
		var length int
		var videoType int

		err = rows.Scan(&id, &queriedTitle, &length, &videoType)
		checkError(err)

		var title [102]uint16
		tempTitle := utf16.Encode([]rune(queriedTitle))
		copy(title[:], tempTitle)

		l.NewVideoTable = append(l.NewVideoTable, NewVideoTable{
			ID:          uint32(id),
			VideoLength: uint16(length),
			TitleID:     0,
			Unknown:     [15]byte{},
			Unknown2:    0,
			RatingID:    9,
			Unknown3:    1,
			Title:       title,
		})
	}

	l.Header.NumberOfNewVideoTables = uint32(len(l.NewVideoTable))
}

func (l *List) MakePopularVideoTable() {
	l.Header.PopularVideoTableOffset = l.GetCurrentSize()

	rows, err := pool.Query(ctx, l.getVideoQueryString())
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		var id int
		var queriedTitle string
		var length int
		var videoType int

		err = rows.Scan(&id, &queriedTitle, &length, &videoType)
		checkError(err)

		var title [102]uint16
		tempTitle := utf16.Encode([]rune(queriedTitle))
		copy(title[:], tempTitle)

		l.PopularVideosTable = append(l.PopularVideosTable, PopularVideosTable{
			ID:          uint32(id),
			VideoLength: uint16(length),
			TitleID:     0,
			BarColor:    0,
			RatingID:    9,
			Unknown:     1,
			VideoRank:   1,
			Unknown2:    222,
			Title:       title,
		})
	}

	l.Header.NumberOfPopularVideoTables = uint32(len(l.PopularVideosTable))
}
