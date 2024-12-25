package info

import "strings"

type SupportedLanguages struct {
	Chinese  uint8
	Korean   uint8
	Japanese uint8
	English  uint8
	French   uint8
	Spanish  uint8
	German   uint8
	Italian  uint8
	Dutch    uint8
}

func (i *Info) GetSupportedLanguages(languages string) {
	for _, s := range strings.Split(languages, ",") {
		switch s {
		case "ZHCN":
			i.SupportedLanguages.Chinese = 1
			break
		case "KO":
			i.SupportedLanguages.Korean = 1
			break
		case "JA":
			i.SupportedLanguages.Japanese = 1
			break
		case "EN":
			i.SupportedLanguages.English = 1
			break
		case "FR":
			i.SupportedLanguages.French = 1
			break
		case "ES":
			i.SupportedLanguages.Spanish = 1
			break
		case "DE":
			i.SupportedLanguages.German = 1
			break
		case "IT":
			i.SupportedLanguages.Italian = 1
			break
		case "NL":
			i.SupportedLanguages.Dutch = 1
			break
		}
	}
}
