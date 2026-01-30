package processing

import (
	"ballot-tool/internal/models"
	"sort"
)

func SortByCloses(rows []models.BallotWithRole) {
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Ballot.Closing.Before(rows[j].Ballot.Closing)
	})
}
