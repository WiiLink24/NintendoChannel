package dllist

import (
	"NintendoChannel/constants"
	"unicode/utf16"
)

// TitleTypeTable contains all the title types the Nintendo Channel can display
type TitleTypeTable struct {
	TypeID       uint8
	ConsoleModel constants.ConsoleModels
	ConsoleName  [51]uint16
	GroupID      constants.TitleGroupTypes
	Unknown      uint8
}

func (l *List) MakeTitleTypeTable() {
	l.Header.TitleTypeTableOffset = l.GetCurrentSize()

	for _, titleType := range constants.TitleTypesData {
		var consoleNameFinal [51]uint16
		consoleName := utf16.Encode([]rune(titleType.ConsoleName))
		copy(consoleNameFinal[:], consoleName)

		data := TitleTypeTable{
			TypeID:       titleType.TypeID,
			ConsoleModel: titleType.ConsoleModel,
			ConsoleName:  consoleNameFinal,
			GroupID:      titleType.GroupID,
			Unknown:      255,
		}

		l.TitleTypesTable = append(l.TitleTypesTable, data)
	}

	l.Header.NumberOfTitleTypeTables = uint32(len(l.TitleTypesTable))
}
