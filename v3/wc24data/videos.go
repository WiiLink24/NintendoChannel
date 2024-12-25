package wc24data

import (
	"NintendoChannel/common"
	"NintendoChannel/constants"
	"time"
	"unicode/utf16"
)

type VideoTable struct {
	ID          uint32
	VideoLength uint16
	TitleID     uint32
	VideoType   uint8
	Unknown     uint8
	IsNotTitle  uint8
	_           uint8
	IsNew       uint8
	Title       [51]uint16
}

// IMPORTANT: Video types differ between v6 NC and v3 NC. A video type of 6 is a message in v3 NC.

func (w *WC24Data) MakeVideoTable() {
	w.Header.VideoTableOffset = w.GetCurrentSize()

	rows, err := pool.Query(ctx, constants.GetVideoQueryString(w.language))
	common.CheckError(err)

	index := 1
	defer rows.Close()
	for rows.Next() {
		var id int
		var queriedTitle string
		var length int
		var videoType int
		var dateAdded time.Time

		err = rows.Scan(&id, &queriedTitle, &length, &videoType, &dateAdded)
		common.CheckError(err)

		var title [51]uint16
		tempTitle := utf16.Encode([]rune(queriedTitle))
		copy(title[:], tempTitle)

		w.VideoTable = append(w.VideoTable, VideoTable{
			ID:          uint32(id),
			VideoLength: uint16(length),
			TitleID:     0,
			VideoType:   uint8(videoType),
			Unknown:     0xFF,
			IsNotTitle:  1,
			IsNew:       0,
			Title:       title,
		})
		index++
	}

	w.Header.NumberOfVideos = uint32(len(w.VideoTable))
}
