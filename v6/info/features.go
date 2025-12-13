package info

import (
	"NintendoChannel/gametdb"
	"slices"
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

func (i *Info) GetSupportedFeatures(game *gametdb.Game) {
	touchGenIDs := []string{
		"YBN", "VAA", "AYA", "AND", "ANM", "ATD", "CVN",
		"YCU", "ATI", "AOS", "AG3", "AWI", "APL", "AJQ", "CM7",
		"AD5", "AD2", "ADG", "AD7", "AD3", "IMW", "C6P", "AXP",
		"A8N", "AZI", "ASQ", "ATR", "AGF",
		"RFN", "RFP", "R64", "RYW",
	}

	for _, s := range game.Features.Feature {
		if strings.Contains(s, "online") {
			i.SupportedFeatures.Online = 1
			i.SupportedFeatures.NintendoWifiConnection = 1
		} else if s == "download" {
			i.SupportedFeatures.DLC = 1
		}
	}

	if game.Controllers.MultiCart > 1 {
		i.SupportedFeatures.WirelessPlay = 1
	}
	if game.Controllers.SingleCart > 1 {
		i.SupportedFeatures.DownloadPlay = 1
	}

	if slices.Contains(touchGenIDs, game.ID[:3]) {
		i.SupportedFeatures.TouchGenerationsTitle = 1
	}
}
