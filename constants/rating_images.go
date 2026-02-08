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

	//go:embed ratings/esrb/descriptor.otf
	ESRBRatingDescriptorFont []byte

	//go:embed ratings/pegi/descriptors/discrimination.jpg
	PEGIDiscrimination []byte

	//go:embed ratings/pegi/descriptors/drugs.jpg
	PEGIDrugs []byte

	//go:embed ratings/pegi/descriptors/fear.jpg
	PEGIFear []byte

	//go:embed ratings/pegi/descriptors/gambling.jpg
	PEGIGambling []byte

	//go:embed ratings/pegi/descriptors/bad-language.jpg
	PEGILanguage []byte

	//go:embed ratings/pegi/descriptors/sexual-content.jpg
	PEGISexualContent []byte

	//go:embed ratings/pegi/descriptors/violence.jpg
	PEGIViolence []byte

	//go:embed ratings/pegi/descriptors/online.jpg
	PEGIOnline []byte

	PEGIDescriptors = map[string][]byte{
		"discrimination": PEGIDiscrimination,
		"drugs":          PEGIDrugs,
		"fear":           PEGIFear,
		"gambling":       PEGIGambling,
		"language":       PEGILanguage,
		"sex":            PEGISexualContent,
		"violence":       PEGIViolence,
		"online":         PEGIOnline,
	}

	//go:embed ratings/cero/descriptors/crime.jpg
	CEROCrime []byte

	//go:embed ratings/cero/descriptors/drinkingandsmoking.jpg
	CERODrinkingAndSmoking []byte

	//go:embed ratings/cero/descriptors/drugs.jpg
	CERODrugs []byte

	//go:embed ratings/cero/descriptors/fright.jpg
	CEROFright []byte

	//go:embed ratings/cero/descriptors/gambling.jpg
	CEROGambling []byte

	//go:embed ratings/cero/descriptors/language.jpg
	CEROLanguage []byte

	//go:embed ratings/cero/descriptors/love.jpg
	CEROLove []byte

	//go:embed ratings/cero/descriptors/sex.jpg
	CEROSex []byte

	//go:embed ratings/cero/descriptors/violence.jpg
	CEROViolence []byte

	CERODescriptors = map[string][]byte{
		"crime":           CEROCrime,
		"alcohol/tobacco": CERODrinkingAndSmoking,
		"drugs":           CERODrugs,
		"horror":          CEROFright,
		"gambling":        CEROGambling,
		"language":        CEROLanguage,
		"love":            CEROLove,
		"sex":             CEROSex,
		"violence":        CEROViolence,
	}

	Images = map[RatingGroup][][]byte{
		CERO: {CEROA, CEROB, CEROC, CEROD, CEROZ},
		ESRB: {ECImage, EImage, E10Image, TImage, MImage},
		PEGI: {PEGI3, PEGI7, PEGI12, PEGI16, PEGI18},
	}
)
