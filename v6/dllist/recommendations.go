package dllist

import (
	"NintendoChannel/common"
	"NintendoChannel/constants"
	"fmt"
)

type RecentRecommendationTable struct {
	TitleOffset uint32
	Medal       constants.Medal
	Unknown     uint8
}

const QueryRecommendations = `SELECT COUNT(*), game_id AS game_count FROM recommendations GROUP BY game_id`

// BaseRecommendationColumnQuery is a query that allows for getting the amount of votes for any table.
// $1: column value
// $2: game id
// $3: gender
// $4: lower age bound
// $5 upper age bound
const _BaseRecommendationColumnQuery = `SELECT COUNT(%s) FROM recommendations WHERE %s = $1
                                        AND game_id = $2
                                        AND gender = $3 
                                        AND age >= $4 
                                        AND age <= $5`

// BaseRecommendationColumnQueryNoGender is a query that allows for getting the amount of votes for any table.
// $1: column value
// $2: game id
// $3: lower age bound
// $4 upper age bound
const _BaseRecommendationColumnQueryNoGender = `SELECT COUNT(%s) FROM recommendations WHERE %s = $1 
                                        AND game_id = $2 
                                        AND age >= $3
                                        AND age <= $4`

func BaseRecommendationColumnQueryNoGender(columnName string) string {
	return fmt.Sprintf(_BaseRecommendationColumnQueryNoGender, columnName, columnName)
}

func BaseRecommendationColumnQuery(columnName string) string {
	return fmt.Sprintf(_BaseRecommendationColumnQuery, columnName, columnName)
}

func PopulateRecommendations() {
	rows, err := pool.Query(ctx, QueryRecommendations)
	common.CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var gameID string
		var count int
		err = rows.Scan(&count, &gameID)
		common.CheckError(err)

		recommendations[gameID] = common.TitleRecommendation{
			NumberOfRecommendations: count,
			AllRecommendations:      constants.AgeRecommendationTable,
			MaleRecommendations:     constants.AgeRecommendationTable,
			FemaleRecommendations:   constants.AgeRecommendationTable,
		}

		// We now have to query for all the different types of recommendation criteria
		// First start with all.
		for i, rec := range recommendations[gameID].AllRecommendations {
			// Appeal
			var everyone int
			err := pool.QueryRow(ctx, BaseRecommendationColumnQueryNoGender("appeal"), 0, gameID, rec.LowerAge, rec.UpperAge).Scan(&everyone)
			if err != nil {
				panic(err)
			}

			var gamers int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQueryNoGender("appeal"), 1, gameID, rec.LowerAge, rec.UpperAge).Scan(&gamers)
			if err != nil {
				panic(err)
			}

			if gamers+everyone != 0 {
				if entry, ok := recommendations[gameID]; ok {
					entry.AllRecommendations[i].EveryonePercent = uint8((float64(everyone) / float64(everyone+gamers)) * 100)
					recommendations[gameID] = entry
				}
			}

			if everyone < gamers {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.AllRecommendations[i].IsGamers = constants.True
					recommendations[gameID] = entry
				}
			} else if everyone != gamers {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.AllRecommendations[i].IsGamers = constants.False
					recommendations[gameID] = entry
				}
			}

			// Mood
			var casual int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQueryNoGender("gaming_mood"), 0, gameID, rec.LowerAge, rec.UpperAge).Scan(&casual)
			if err != nil {
				panic(err)
			}

			var hardcore int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQueryNoGender("gaming_mood"), 1, gameID, rec.LowerAge, rec.UpperAge).Scan(&hardcore)
			if err != nil {
				panic(err)
			}

			if casual+hardcore != 0 {
				if entry, ok := recommendations[gameID]; ok {
					entry.AllRecommendations[i].CasualPercent = uint8((float64(casual) / float64(casual+hardcore)) * 100)
					recommendations[gameID] = entry
				}
			}

			if casual < hardcore {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.AllRecommendations[i].IsHardcore = constants.True
					recommendations[gameID] = entry
				}
			} else if casual != hardcore {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.AllRecommendations[i].IsHardcore = constants.False
					recommendations[gameID] = entry
				}
			}

			// Friend or Alone
			var alone int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQueryNoGender("friend_or_alone"), 0, gameID, rec.LowerAge, rec.UpperAge).Scan(&alone)
			if err != nil {
				panic(err)
			}

			var friend int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQueryNoGender("friend_or_alone"), 1, gameID, rec.LowerAge, rec.UpperAge).Scan(&friend)
			if err != nil {
				panic(err)
			}

			if alone+friend != 0 {
				if entry, ok := recommendations[gameID]; ok {
					entry.AllRecommendations[i].AlonePercent = uint8((float64(alone) / float64(alone+friend)) * 100)
					recommendations[gameID] = entry
				}
			}

			if alone < friend {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.AllRecommendations[i].IsWithFriends = constants.True
					recommendations[gameID] = entry
				}
			} else if alone != friend {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.AllRecommendations[i].IsWithFriends = constants.False
					recommendations[gameID] = entry
				}
			}
		}

		// Next is Male.
		for i, rec := range recommendations[gameID].MaleRecommendations {
			// Appeal
			var everyone int
			err := pool.QueryRow(ctx, BaseRecommendationColumnQuery("appeal"), 0, gameID, 1, rec.LowerAge, rec.UpperAge).Scan(&everyone)
			if err != nil {
				panic(err)
			}

			var gamers int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQuery("appeal"), 1, gameID, 1, rec.LowerAge, rec.UpperAge).Scan(&gamers)
			if err != nil {
				panic(err)
			}

			if gamers+everyone != 0 {
				if entry, ok := recommendations[gameID]; ok {
					entry.MaleRecommendations[i].EveryonePercent = uint8((float64(everyone) / float64(everyone+gamers)) * 100)
					recommendations[gameID] = entry
				}
			}

			if everyone < gamers {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.MaleRecommendations[i].IsGamers = constants.True
					recommendations[gameID] = entry
				}
			} else if everyone != gamers {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.MaleRecommendations[i].IsGamers = constants.False
					recommendations[gameID] = entry
				}
			}

			// Mood
			var casual int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQuery("gaming_mood"), 0, gameID, 1, rec.LowerAge, rec.UpperAge).Scan(&casual)
			if err != nil {
				panic(err)
			}

			var hardcore int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQuery("gaming_mood"), 1, gameID, 1, rec.LowerAge, rec.UpperAge).Scan(&hardcore)
			if err != nil {
				panic(err)
			}

			if casual+hardcore != 0 {
				if entry, ok := recommendations[gameID]; ok {
					entry.MaleRecommendations[i].CasualPercent = uint8((float64(casual) / float64(casual+hardcore)) * 100)
					recommendations[gameID] = entry
				}
			}

			if casual < hardcore {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.MaleRecommendations[i].IsHardcore = constants.True
					recommendations[gameID] = entry
				}
			} else if casual != hardcore {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.MaleRecommendations[i].IsHardcore = constants.False
					recommendations[gameID] = entry
				}
			}

			// Friend or Alone
			var alone int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQuery("friend_or_alone"), 0, gameID, 1, rec.LowerAge, rec.UpperAge).Scan(&alone)
			if err != nil {
				panic(err)
			}

			var friend int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQuery("friend_or_alone"), 1, gameID, 1, rec.LowerAge, rec.UpperAge).Scan(&friend)
			if err != nil {
				panic(err)
			}

			if alone+friend != 0 {
				if entry, ok := recommendations[gameID]; ok {
					entry.MaleRecommendations[i].AlonePercent = uint8((float64(alone) / float64(alone+friend)) * 100)
					recommendations[gameID] = entry
				}
			}

			if alone < friend {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.MaleRecommendations[i].IsWithFriends = constants.True
					recommendations[gameID] = entry
				}
			} else if alone != friend {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.MaleRecommendations[i].IsWithFriends = constants.False
					recommendations[gameID] = entry
				}
			}
		}

		// Finally is Female.
		for i, rec := range recommendations[gameID].FemaleRecommendations {
			// Appeal
			var everyone int
			err := pool.QueryRow(ctx, BaseRecommendationColumnQuery("appeal"), 0, gameID, 2, rec.LowerAge, rec.UpperAge).Scan(&everyone)
			if err != nil {
				panic(err)
			}

			var gamers int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQuery("appeal"), 1, gameID, 2, rec.LowerAge, rec.UpperAge).Scan(&gamers)
			if err != nil {
				panic(err)
			}

			if gamers+everyone != 0 {
				if entry, ok := recommendations[gameID]; ok {
					entry.FemaleRecommendations[i].EveryonePercent = uint8((float64(everyone) / float64(everyone+gamers)) * 100)
					recommendations[gameID] = entry
				}
			}

			if everyone < gamers {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.FemaleRecommendations[i].IsGamers = constants.True
					recommendations[gameID] = entry
				}
			} else if everyone != gamers {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.FemaleRecommendations[i].IsGamers = constants.False
					recommendations[gameID] = entry
				}
			}

			// Mood
			var casual int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQuery("gaming_mood"), 0, gameID, 2, rec.LowerAge, rec.UpperAge).Scan(&casual)
			if err != nil {
				panic(err)
			}

			var hardcore int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQuery("gaming_mood"), 1, gameID, 2, rec.LowerAge, rec.UpperAge).Scan(&hardcore)
			if err != nil {
				panic(err)
			}

			if casual+hardcore != 0 {
				if entry, ok := recommendations[gameID]; ok {
					entry.FemaleRecommendations[i].CasualPercent = uint8((float64(casual) / float64(casual+hardcore)) * 100)
					recommendations[gameID] = entry
				}
			}

			if casual < hardcore {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.FemaleRecommendations[i].IsHardcore = constants.True
					recommendations[gameID] = entry
				}
			} else if casual != hardcore {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.FemaleRecommendations[i].IsHardcore = constants.False
					recommendations[gameID] = entry
				}
			}

			// Friend or Alone
			var alone int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQuery("friend_or_alone"), 0, gameID, 2, rec.LowerAge, rec.UpperAge).Scan(&alone)
			if err != nil {
				panic(err)
			}

			var friend int
			err = pool.QueryRow(ctx, BaseRecommendationColumnQuery("friend_or_alone"), 1, gameID, 2, rec.LowerAge, rec.UpperAge).Scan(&friend)
			if err != nil {
				panic(err)
			}

			if alone+friend != 0 {
				if entry, ok := recommendations[gameID]; ok {
					entry.FemaleRecommendations[i].AlonePercent = uint8((float64(alone) / float64(alone+friend)) * 100)
					recommendations[gameID] = entry
				}
			}

			if alone < friend {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.FemaleRecommendations[i].IsWithFriends = constants.True
					recommendations[gameID] = entry
				}
			} else if alone != friend {
				if entry, ok := recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.FemaleRecommendations[i].IsWithFriends = constants.False
					recommendations[gameID] = entry
				}
			}
		}
	}
}

func (l *List) MakeRecommendationTable() {
	l.Header.RecommendationTableOffset = l.GetCurrentSize()

	for gameID, _ := range recommendations {
		// Now we find the title from our title table
		for i, title := range l.TitleTable {
			if string(title.TitleID[:]) == gameID {
				l.RecommendationTable = append(l.RecommendationTable, (236*uint32(i))+l.Header.TitleTableOffset)
				break
			}
		}
	}

	l.Header.NumberOfRecommendationTables = uint32(len(l.RecommendationTable))
}

func (l *List) MakeRecentRecommendationTable() {
	l.Header.RecentRecommendationTableOffset = l.GetCurrentSize()

	for gameID, rec := range recommendations {
		for i, title := range l.TitleTable {
			if string(title.TitleID[:]) == gameID {
				l.RecentRecommendationTable = append(l.RecentRecommendationTable, RecentRecommendationTable{
					TitleOffset: (236 * uint32(i)) + l.Header.TitleTableOffset,
					Medal:       GetMedal(rec, rec.NumberOfRecommendations),
					Unknown:     222,
				})
				break
			}
		}
	}

	l.Header.NumberOfRecentRecommendationTables = uint32(len(l.RecentRecommendationTable))
}
