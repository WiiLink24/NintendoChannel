package info

import (
	"NintendoChannel/common"
	"NintendoChannel/constants"
)

type RecommendationTable struct {
	// [3]Gender[8]RecommendationPercent
	Everyone [3][8]uint8
	Casual   [3][8]uint8
	Alone    [3][8]uint8
	Medals   [3][8]uint8
}

func (i *Info) MakeRecommendationTable(recommendation common.TitleRecommendation, numberOfTimesVotes int) {
	var medal constants.Medal

	if numberOfTimesVotes <= 20 { // Don't even bother calculating
		return
	}

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

			avg := (int(i.RecommendationTable.Everyone[j][k]) + int(i.RecommendationTable.Casual[j][k]) + int(i.RecommendationTable.Alone[j][k])) / 3
			switch {
			case avg >= 90:
				medal = constants.Platinum
			case avg >= 80:
				medal = constants.Gold
			case avg >= 70:
				medal = constants.Silver
			case avg >= 60:
				medal = constants.Bronze
			default:
				medal = constants.None
			}
			i.RecommendationTable.Medals[j][k] = uint8(medal)
		}
	}
}
