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

func (i *Info) GetSupportedFeatures(features *gametdb.Features) {
	for _, s := range features.Feature {
		if strings.Contains(s, "online") {
			i.SupportedFeatures.Online = 1
			i.SupportedFeatures.NintendoWifiConnection = 1
		} else if s == "nintendods" {
			i.SupportedFeatures.DownloadPlay = 1
		} else if s == "download" {
			i.SupportedFeatures.DLC = 1
		}
	}
}
