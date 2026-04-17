package ballot

import (
	"sort"
)

func SortByCloses(rows []BallotWithRole) {
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Ballot.Closing.Before(rows[j].Ballot.Closing)
	})
}
