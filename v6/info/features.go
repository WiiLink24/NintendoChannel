package info

import (
	"NintendoChannel/constants"
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

	for _, s := range game.Features.Feature {
		if strings.Contains(s, "online") {
			i.SupportedFeatures.Online = 1
			i.SupportedFeatures.NintendoWifiConnection = 1
		} else if s == "download" {
			i.SupportedFeatures.DLC = 1
		}
	}

	if slices.Contains(constants.PaynPlayIDs, game.ID[:3]) {
		i.SupportedFeatures.NintendoWifiConnection = 2
	}

	if game.Controllers.MultiCart > 1 {
		i.SupportedFeatures.WirelessPlay = 1
	}
	if game.Controllers.SingleCart > 1 {
		i.SupportedFeatures.DownloadPlay = 1
	}

	if slices.Contains(constants.TouchGenIDs, game.ID[:3]) {
		i.SupportedFeatures.TouchGenerationsTitle = 1
	}
}
