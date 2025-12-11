package info

import (
	"NintendoChannel/gametdb"
	"strings"
)

type SupportedFeatures struct {
	Miis                   uint8
	Online                 uint8
	WiiConnect24           uint8
	NintendoWifiConnection uint8
	DLC                    uint8
	WirelessPlay           uint8
	DownloadPlay           uint8
	TouchGenerationsTitle  uint8
}

func (i *Info) GetSupportedFeatures(features *gametdb.Features, gameID string, controllers *gametdb.Controllers) {
	touchGenIDs := []string{
		"YBN", "VAA", "AYA", "AND", "ANM", "ATD", "CVN",
		"YCU", "ATI", "AOS", "AG3", "AWI", "APL", "AJQ", "CM7",
		"AD5", "AD2", "ADG", "AD7", "AD3", "IMW", "C6P", "AXP",
		"A8N", "AZI", "ASQ", "ATR", "AGF",
		"RFN", "RFP", "R64", "RYW",
	}

	for _, s := range features.Feature {
		if strings.Contains(s, "online") {
			i.SupportedFeatures.Online = 1
			i.SupportedFeatures.NintendoWifiConnection = 1
		} else if s == "download" {
			i.SupportedFeatures.DLC = 1
		}
	}

	if controllers.MultiCart >= 2 && controllers.MultiCart <= 8 {
		i.SupportedFeatures.WirelessPlay = 1
	}
	if controllers.SingleCart >= 1 && controllers.SingleCart <= 8 {
		i.SupportedFeatures.DownloadPlay = 1
	}

	prefix := ""
	if len(gameID) >= 3 {
		prefix = gameID[:3]
	}

	for _, id := range touchGenIDs {
		if prefix == id {
			i.SupportedFeatures.TouchGenerationsTitle = 1
			break
		}
	}
}
