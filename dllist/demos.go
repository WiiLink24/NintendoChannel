package dllist

import (
	"NintendoChannel/common"
	"NintendoChannel/constants"
	"bytes"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	"unicode/utf16"
)

type DemoTable struct {
	ID            uint32
	Title         [31]uint16
	Subtitle      [31]uint16
	TitleID       uint32
	CompanyOffset uint32
	RemovalYear   uint16
	RemovalMonth  uint8
	RemovalDay    uint8
	_             uint32
	RatingID      uint8
	IsNew         uint8
	_             uint8
	_             [205]byte
}

var demosTable []*NinchDllist_DemosTable

func PopulateDemos() {
	dl := NewNinchDllist()
	err := dl.Read(kaitai.NewStream(bytes.NewReader(constants.DLList)), nil, dl)
	common.CheckError(err)

	demosTable, err = dl.DemosTable()
	common.CheckError(err)
}

func (l *List) MakeDemoTable() {
	l.Header.DemoTableOffset = l.GetCurrentSize()

	for _, demo := range demosTable {
		var title [31]uint16
		tempTitle := utf16.Encode([]rune(demo.Title))
		copy(title[:], tempTitle)

		var subtitle [31]uint16
		tempSubtitle := utf16.Encode([]rune(demo.Subtitle))
		copy(subtitle[:], tempSubtitle)

		l.DemoTable = append(l.DemoTable, DemoTable{
			ID:            demo.Id,
			Title:         title,
			Subtitle:      subtitle,
			TitleID:       demo.Titleid,
			CompanyOffset: l.Header.CompanyTableOffset,
			RemovalYear:   0xFFFF,
			RemovalMonth:  0xFF,
			RemovalDay:    0xFF,
			RatingID:      demo.RatingId,
			IsNew:         0,
		})
	}

	l.Header.NumberOfDemoTables = uint32(len(l.DemoTable))
}
