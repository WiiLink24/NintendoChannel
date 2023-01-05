package dllist

import (
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

func (l *List) MakeDemoTable() {
	l.Header.DemoTableOffset = l.GetCurrentSize()

	var title [31]uint16
	tempTitle := utf16.Encode([]rune("The Ghost"))
	copy(title[:], tempTitle)

	l.DemoTable = append(l.DemoTable, DemoTable{
		ID:            1,
		Title:         title,
		Subtitle:      [31]uint16{},
		TitleID:       0,
		CompanyOffset: l.Header.CompanyTableOffset,
		RemovalYear:   0xFFFF,
		RemovalMonth:  0xFF,
		RemovalDay:    0xFF,
		RatingID:      9,
		IsNew:         0,
	})

	l.Header.NumberOfDemoTables = 1
}
