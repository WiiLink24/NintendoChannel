package info

import (
	"NintendoChannel/gametdb"
	"strings"
	"unicode/utf16"
)

type SupportedControllers struct {
	WiiRemote          uint8
	Nunchuk            uint8
	ClassicController  uint8
	GamecubeController uint8
}

func (i *Info) GetSupportedControllers(controllers *gametdb.Controllers) {
	wrotePeripheral := []string{}
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
		// All of these are set to append in case there's multiple peripherals
		case "wheel":
			wrotePeripheral = append(wrotePeripheral, "Wii Wheel")
		case "balanceboard":
			wrotePeripheral = append(wrotePeripheral, "Wii Balance Board")
		case "wiispeak":
			wrotePeripheral = append(wrotePeripheral, "Wii Speak")
		case "zapper":
			wrotePeripheral = append(wrotePeripheral, "Wii Zapper")
		case "dancepad":
			wrotePeripheral = append(wrotePeripheral, "Dance Pad")
		case "guitar":
			wrotePeripheral = append(wrotePeripheral, "Guitar")
		case "drums":
			wrotePeripheral = append(wrotePeripheral, "Drums")
		case "microphone":
			wrotePeripheral = append(wrotePeripheral, "Microphone")
		case "keyboard":
			wrotePeripheral = append(wrotePeripheral, "USB Keyboard")
		}
	}
	if len(wrotePeripheral) > 0 {
		peripheralsString := strings.Join(wrotePeripheral, ", ")
		// For some reason, the peripheral text must be padded with 2 uint16 before any real text
		temp := []uint16{0, 0}
		encodedText := append(temp, utf16.Encode([]rune(peripheralsString))...)
		copy(i.PeripheralsText[:], encodedText)
	}
}
