package common

import "NintendoChannel/constants"

type TitleRecommendation struct {
	NumberOfRecommendations int
	AllRecommendations      [8]constants.AgeRecommendationData
	MaleRecommendations     [8]constants.AgeRecommendationData
	FemaleRecommendations   [8]constants.AgeRecommendationData
}
