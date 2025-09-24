package info

import (
	"NintendoChannel/common"
)

type RecommendationTable struct {
	// [3]Gender[8]RecommendationPercent
	Everyone [3][8]uint8
	Casual   [3][8]uint8
	Alone    [3][8]uint8
	Medals   [3][8]uint8
}

func (i *Info) MakeRecommendationTable(recommendation common.TitleRecommendation) {
	for j := 0; j < 3; j++ {
		for k := 0; k < 8; k++ {
			if j == 0 {
				i.RecommendationTable.Everyone[j][k] = recommendation.AllRecommendations[k].EveryonePercent
				i.RecommendationTable.Casual[j][k] = recommendation.AllRecommendations[k].CasualPercent
				i.RecommendationTable.Alone[j][k] = recommendation.AllRecommendations[k].AlonePercent
			} else if j == 1 {
				i.RecommendationTable.Everyone[j][k] = recommendation.MaleRecommendations[k].EveryonePercent
				i.RecommendationTable.Casual[j][k] = recommendation.MaleRecommendations[k].CasualPercent
				i.RecommendationTable.Alone[j][k] = recommendation.MaleRecommendations[k].AlonePercent
			} else {
				i.RecommendationTable.Everyone[j][k] = recommendation.FemaleRecommendations[k].EveryonePercent
				i.RecommendationTable.Casual[j][k] = recommendation.FemaleRecommendations[k].CasualPercent
				i.RecommendationTable.Alone[j][k] = recommendation.FemaleRecommendations[k].AlonePercent
			}
		}
	}
}
