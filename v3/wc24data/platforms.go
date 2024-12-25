package wc24data

import (
	"NintendoChannel/constants"
	"unicode/utf16"
)

type Platform struct {
	TypeID       uint8
	Unk          uint8
	PlatformName [29]uint16
	GroupId      uint32
}

func (w *WC24Data) MakePlatformTable() {
	w.Header.PlatformTableOffset = w.GetCurrentSize()

	for _, titleType := range constants.TitleTypesData {
		var consoleNameFinal [29]uint16
		consoleName := utf16.Encode([]rune(titleType.ConsoleName))
		copy(consoleNameFinal[:], consoleName)

		data := Platform{
			TypeID:       titleType.TypeID,
			Unk:          48,
			PlatformName: consoleNameFinal,
			GroupId:      uint32(titleType.GroupID),
		}

		w.PlatformTable = append(w.PlatformTable, data)
	}

	w.Header.NumberOfPlatforms = uint32(len(w.PlatformTable))
}
