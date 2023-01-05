package dllist

type RecentRecommendationTable struct {
	TitleOffset uint32
	Medal       uint8
	Unknown     uint8
}

func (l *List) MakeRecommendationTable() {
	l.Header.RecommendationTableOffset = l.GetCurrentSize()
	l.RecommendationTable = append(l.RecommendationTable, l.Header.TitleTableOffset)
	l.Header.NumberOfRecommendationTables = 1
}

func (l *List) MakeRecentRecommendationTable() {
	l.Header.RecentRecommendationTableOffset = l.GetCurrentSize()
	l.RecentRecommendationTable = append(l.RecentRecommendationTable, RecentRecommendationTable{
		TitleOffset: l.Header.TitleTableOffset,
		Medal:       0,
		Unknown:     222,
	})
	l.Header.NumberOfRecentRecommendationTables = 1
}
