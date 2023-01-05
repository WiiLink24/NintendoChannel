package constants

import _ "embed"

var (
	//go:embed nc.bin
	DLList []byte

	//go:embed ratings/esrb/EC.jpg
	ECImage []byte

	//go:embed ratings/esrb/E.jpg
	EImage []byte

	//go:embed ratings/esrb/E10.jpg
	E10Image []byte

	//go:embed ratings/esrb/T.jpg
	TImage []byte

	//go:embed ratings/esrb/M.jpg
	MImage []byte

	//go:embed ratings/pegi/3.jpg
	PEGI3 []byte

	//go:embed ratings/pegi/7.jpg
	PEGI7 []byte

	//go:embed ratings/pegi/12.jpg
	PEGI12 []byte

	//go:embed ratings/pegi/16.jpg
	PEGI16 []byte

	//go:embed ratings/pegi/18.jpg
	PEGI18 []byte

	//go:embed ratings/cero/A.jpg
	CEROA []byte

	//go:embed ratings/cero/B.jpg
	CEROB []byte

	//go:embed ratings/cero/C.jpg
	CEROC []byte

	//go:embed ratings/cero/D.jpg
	CEROD []byte

	//go:embed ratings/cero/Z.jpg
	CEROZ []byte

	Images = map[RatingGroup][][]byte{
		CERO: {CEROA, CEROB, CEROC, CEROD, CEROZ},
		ESRB: {ECImage, EImage, E10Image, TImage, MImage},
		PEGI: {PEGI3, PEGI7, PEGI12, PEGI16, PEGI18},
	}
)
