package processing

import (
	"ballot-tool/internal/models"
	"log"
)

func JoinBallotRole(roles []models.Role, ballots []models.Ballot) ([]models.BallotWithRole, []models.Ballot) {
	roleCommitteeIdx := make(map[string]*models.Role, len(roles))
	for i := range roles {
		c := roles[i].Committee
		if c == "" {
			continue
		}
		if _, exists := roleCommitteeIdx[c]; !exists {
			roleCommitteeIdx[c] = &roles[i]
		}
	}

	matches := make([]models.BallotWithRole, 0, len(ballots))
	missing := make([]models.Ballot, 0, len(ballots))
	for _, b := range ballots {
		match, ok := roleCommitteeIdx[b.Committee]
		if !ok {
			missing = append(missing, b)
			continue
		}

		matches = append(matches, models.BallotWithRole{
			Ballot: b,
			Role:   *match,
		})
	}

	log.Printf("matched %d ballots \n", len(matches))
	log.Printf("found %d ballots without voter \n", len(missing))

	return matches, missing
}
