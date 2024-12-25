package info

import (
	"NintendoChannel/gametdb"
	"unicode/utf16"
)

type SupportedControllers struct {
	WiiRemote          uint8
	Nunchuk            uint8
	ClassicController  uint8
	GamecubeController uint8
}

func (i *Info) GetSupportedControllers(controllers *gametdb.Controllers) {
	wrotePeripheral := false
	for _, s := range controllers.Controller {
		switch s.Type {
		case "wiimote":
			i.SupportedControllers.WiiRemote = 1
			break
		case "nunchuk":
			i.SupportedControllers.Nunchuk = 1
			break
		case "classiccontroller":
			i.SupportedControllers.ClassicController = 1
			break
		case "gamecube":
			i.SupportedControllers.GamecubeController = 1
			break
		case "mii":
			// Mii's aren't a controller, but they are considered one to GameTDB for some reason
			i.SupportedFeatures.Miis = 1
			break
		case "wheel":
			if !wrotePeripheral {
				// For some reason the peripheral text must be padded with 2 uint16 before any real text.
				temp := []uint16{0, 0}
				copy(i.PeripheralsText[:], append(temp, utf16.Encode([]rune("Wii Wheel"))...))
				wrotePeripheral = true
			}
			break
		case "balanceboard":
			if !wrotePeripheral {
				temp := []uint16{0, 0}
				copy(i.PeripheralsText[:], append(temp, utf16.Encode([]rune("Wii Balance Board"))...))
				wrotePeripheral = true
			}
			break
		}
	}
}
