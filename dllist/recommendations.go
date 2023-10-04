package dllist

import (
	"NintendoChannel/constants"
	"fmt"
)

type RecentRecommendationTable struct {
	TitleOffset uint32
	Medal       constants.Medal
	Unknown     uint8
}

type TitleRecommendation struct {
	NumberOfRecommendations int
	AllRecommendations      [8]constants.AgeRecommendationData
	MaleRecommendations     [8]constants.AgeRecommendationData
	FemaleRecommendations   [8]constants.AgeRecommendationData
}

const QueryRecommendations = `SELECT COUNT(game_id), game_id FROM recommendations GROUP BY game_id`

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

func (l *List) QueryRecommendations() {
	rows, err := pool.Query(ctx, QueryRecommendations)
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		var gameID string
		var count int
		err = rows.Scan(&count, &gameID)
		checkError(err)

		// First see if this game could exist in all regions
		isForRegion := false
		if gameID[3:] == "A" || gameID[3:] == "B" || gameID[3:] == "U" || gameID[3:] == "X" {
			isForRegion = true
		} else {
			// Now determine if the game exists for this region
			switch l.region {
			case constants.NTSC:
				if gameID[3:] == "E" || gameID[3:] == "N" {
					isForRegion = true
				}
				break
			case constants.Japan:
				if gameID[3:] == "J" {
					isForRegion = true
				}
				break
			case constants.PAL:
				if gameID[3:] == "P" || gameID[3:] == "L" || gameID[3:] == "M" {
					isForRegion = true
				}
				break
			}
		}

		if !isForRegion {
			continue
		}

		l.recommendations[gameID] = TitleRecommendation{
			NumberOfRecommendations: count,
			AllRecommendations:      constants.AgeRecommendationTable,
			MaleRecommendations:     constants.AgeRecommendationTable,
			FemaleRecommendations:   constants.AgeRecommendationTable,
		}

		// We now have to query for all the different types of recommendation criteria
		// First start with all.
		for i, rec := range l.recommendations[gameID].AllRecommendations {
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

			if everyone < gamers {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.AllRecommendations[i].IsGamers = constants.True
					l.recommendations[gameID] = entry
				}
			} else if everyone != gamers {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.AllRecommendations[i].IsGamers = constants.False
					l.recommendations[gameID] = entry
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

			if casual < hardcore {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.AllRecommendations[i].IsHardcore = constants.True
					l.recommendations[gameID] = entry
				}
			} else if casual != hardcore {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.AllRecommendations[i].IsHardcore = constants.False
					l.recommendations[gameID] = entry
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

			if alone < friend {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.AllRecommendations[i].IsWithFriends = constants.True
					l.recommendations[gameID] = entry
				}
			} else if alone != friend {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.AllRecommendations[i].IsWithFriends = constants.False
					l.recommendations[gameID] = entry
				}
			}
		}

		// Next is Male.
		for i, rec := range l.recommendations[gameID].MaleRecommendations {
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

			if everyone < gamers {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.MaleRecommendations[i].IsGamers = constants.True
					l.recommendations[gameID] = entry
				}
			} else if everyone != gamers {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.MaleRecommendations[i].IsGamers = constants.False
					l.recommendations[gameID] = entry
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

			if casual < hardcore {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.MaleRecommendations[i].IsHardcore = constants.True
					l.recommendations[gameID] = entry
				}
			} else if casual != hardcore {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.MaleRecommendations[i].IsHardcore = constants.False
					l.recommendations[gameID] = entry
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

			if alone < friend {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.MaleRecommendations[i].IsWithFriends = constants.True
					l.recommendations[gameID] = entry
				}
			} else if alone != friend {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.MaleRecommendations[i].IsWithFriends = constants.False
					l.recommendations[gameID] = entry
				}
			}
		}

		// Finally is Female.
		for i, rec := range l.recommendations[gameID].FemaleRecommendations {
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

			if everyone < gamers {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.FemaleRecommendations[i].IsGamers = constants.True
					l.recommendations[gameID] = entry
				}
			} else if everyone != gamers {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.FemaleRecommendations[i].IsGamers = constants.False
					l.recommendations[gameID] = entry
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

			if casual < hardcore {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.FemaleRecommendations[i].IsHardcore = constants.True
					l.recommendations[gameID] = entry
				}
			} else if casual != hardcore {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.FemaleRecommendations[i].IsHardcore = constants.False
					l.recommendations[gameID] = entry
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

			if alone < friend {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.FemaleRecommendations[i].IsWithFriends = constants.True
					l.recommendations[gameID] = entry
				}
			} else if alone != friend {
				if entry, ok := l.recommendations[gameID]; ok {
					// Go does not allow for changing values inside a map.
					entry.FemaleRecommendations[i].IsWithFriends = constants.False
					l.recommendations[gameID] = entry
				}
			}
		}
	}
}

func (l *List) MakeRecommendationTable() {
	l.Header.RecommendationTableOffset = l.GetCurrentSize()

	for gameID, _ := range l.recommendations {
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

	for gameID, num := range l.recommendations {
		for i, title := range l.TitleTable {
			if string(title.TitleID[:]) == gameID {
				l.RecentRecommendationTable = append(l.RecentRecommendationTable, RecentRecommendationTable{
					TitleOffset: (236 * uint32(i)) + l.Header.TitleTableOffset,
					Medal:       GetMedal(num.NumberOfRecommendations),
					Unknown:     222,
				})
				break
			}
		}
	}

	l.Header.NumberOfRecentRecommendationTables = uint32(len(l.RecentRecommendationTable))
}
