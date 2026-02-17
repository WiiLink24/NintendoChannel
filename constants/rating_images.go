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

	//go:embed ratings/esrb/descriptors/alcohol_reference.jpg
	ESRBAlcoholReference []byte

	//go:embed ratings/esrb/descriptors/animated_blood.jpg
	ESRBAnimatedBlood []byte

	//go:embed ratings/esrb/descriptors/animated_violence.jpg
	ESRBAnimatedViolence []byte

	//go:embed ratings/esrb/descriptors/blood_and_gore.jpg
	ESRBBloodandGore []byte

	//go:embed ratings/esrb/descriptors/blood.jpg
	ESRBBlood []byte

	//go:embed ratings/esrb/descriptors/cartoon_violence.jpg
	ESRBCartoonViolence []byte

	//go:embed ratings/esrb/descriptors/comic_mischief.jpg
	ESRBComicMischief []byte

	//go:embed ratings/esrb/descriptors/crude_humor.jpg
	ESRBCrudeHumor []byte

	//go:embed ratings/esrb/descriptors/drug_reference.jpg
	ESRBDrugReference []byte

	//go:embed ratings/esrb/descriptors/edutainment.jpg
	ESRBEdutainment []byte

	//go:embed ratings/esrb/descriptors/fantasy_violence.jpg
	ESRBFantasyViolence []byte

	//go:embed ratings/esrb/descriptors/intense_violence.jpg
	ESRBIntenseViolence []byte

	//go:embed ratings/esrb/descriptors/language.jpg
	ESRBLanguage []byte

	//go:embed ratings/esrb/descriptors/lyrics.jpg
	ESRBLyrics []byte

	//go:embed ratings/esrb/descriptors/mature_humor.jpg
	ESRBMatureHumor []byte

	//go:embed ratings/esrb/descriptors/mild_animated_violence.jpg
	ESRBMildAnimatedViolence []byte

	//go:embed ratings/esrb/descriptors/mild_blood.jpg
	ESRBMildBlood []byte

	//go:embed ratings/esrb/descriptors/mild_cartoon_violence.jpg
	ESRBMildCartoonViolence []byte

	//go:embed ratings/esrb/descriptors/mild_fantasy_violence.jpg
	ESRBMildFantasyViolence []byte

	//go:embed ratings/esrb/descriptors/mild_language.jpg
	ESRBMildLanguage []byte

	//go:embed ratings/esrb/descriptors/mild_lyrics.jpg
	ESRBMildLyrics []byte

	//go:embed ratings/esrb/descriptors/mild_sexual_themes.jpg
	ESRBMildSexualThemes []byte

	//go:embed ratings/esrb/descriptors/mild_suggestive_themes.jpg
	ESRBMildSuggestiveThemes []byte

	//go:embed ratings/esrb/descriptors/mild_violence.jpg
	ESRBMildViolence []byte

	//go:embed ratings/esrb/descriptors/partial_nudity.jpg
	ESRBPartialNudity []byte

	//go:embed ratings/esrb/descriptors/sexual_content.jpg
	ESRBSexualContent []byte

	//go:embed ratings/esrb/descriptors/sexual_themes.jpg
	ESRBSexualThemes []byte

	//go:embed ratings/esrb/descriptors/simulated_gambling.jpg
	ESRBSimulatedGambling []byte

	//go:embed ratings/esrb/descriptors/some_adult_assistance_may_be_needed.jpg
	ESRBSomeAdultAssistanceMayBeNeeded []byte

	//go:embed ratings/esrb/descriptors/strong_language.jpg
	ESRBStrongLanguage []byte

	//go:embed ratings/esrb/descriptors/strong_lyrics.jpg
	ESRBStrongLyrics []byte

	//go:embed ratings/esrb/descriptors/strong_sexual_content.jpg
	ESRBStrongSexualContent []byte

	//go:embed ratings/esrb/descriptors/suggestive_themes.jpg
	ESRBSuggestiveThemes []byte

	//go:embed ratings/esrb/descriptors/tobacco_reference.jpg
	ESRBTobaccoReference []byte

	//go:embed ratings/esrb/descriptors/use_of_alcohol.jpg
	ESRBUseofAlcohol []byte

	//go:embed ratings/esrb/descriptors/use_of_drugs.jpg
	ESRBUseofDrugs []byte

	//go:embed ratings/esrb/descriptors/use_of_tobacco.jpg
	ESRBUseofTobacco []byte

	//go:embed ratings/esrb/descriptors/violence.jpg
	ESRBViolence []byte

	//go:embed ratings/esrb/descriptors/violent_references.jpg
	ESRBViolentReferences []byte

	ESRBDescriptors = map[string][]byte{
		"alcohol reference":                   ESRBAlcoholReference,
		"animated blood":                      ESRBAnimatedBlood,
		"animated violence":                   ESRBAnimatedViolence,
		"blood and gore":                      ESRBBloodandGore,
		"blood":                               ESRBBlood,
		"cartoon violence":                    ESRBCartoonViolence,
		"comic mischief":                      ESRBComicMischief,
		"crude humor":                         ESRBCrudeHumor,
		"drug reference":                      ESRBDrugReference,
		"edutainment":                         ESRBEdutainment,
		"fantasy violence":                    ESRBFantasyViolence,
		"intense violence":                    ESRBIntenseViolence,
		"language":                            ESRBLanguage,
		"lyrics":                              ESRBLyrics,
		"mature humor":                        ESRBMatureHumor,
		"mild animated violence":              ESRBMildAnimatedViolence,
		"mild blood":                          ESRBMildBlood,
		"mild cartoon violence":               ESRBMildCartoonViolence,
		"mild fantasy violence":               ESRBMildFantasyViolence,
		"mild language":                       ESRBMildLanguage,
		"mild lyrics":                         ESRBMildLyrics,
		"mild sexual themes":                  ESRBMildSexualThemes,
		"mild suggestive themes":              ESRBMildSuggestiveThemes,
		"mild violence":                       ESRBMildViolence,
		"partial nudity":                      ESRBPartialNudity,
		"sexual content":                      ESRBSexualContent,
		"sexual themes":                       ESRBSexualThemes,
		"simulated gambling":                  ESRBSimulatedGambling,
		"some adult assistance may be needed": ESRBSomeAdultAssistanceMayBeNeeded,
		"strong language":                     ESRBStrongLanguage,
		"strong lyrics":                       ESRBStrongLyrics,
		"strong sexual content":               ESRBStrongSexualContent,
		"suggestive themes":                   ESRBSuggestiveThemes,
		"tobacco reference":                   ESRBTobaccoReference,
		"use of alcohol":                      ESRBUseofAlcohol,
		"use of drugs":                        ESRBUseofDrugs,
		"use of tobacco":                      ESRBUseofTobacco,
		"violence":                            ESRBViolence,
		"violent references":                  ESRBViolentReferences,
	}
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
