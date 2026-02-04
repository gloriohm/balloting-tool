package processing

import (
	"ballot-tool/internal/models"
	"log"
)

func JoinBallotRole(roles []models.Role, ballots []models.Ballot) ([]models.BallotWithRole, []models.Ballot) {
	roleCommitteeIdx := createRoleComIdx(roles)

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

func JoinCommitteeRole(roles []models.Role, coms []models.Committee) []models.Committee {
	roleCommitteeIdx := createRoleComIdx(roles)

	missing := make([]models.Committee, 0, len(coms))

	for _, c := range coms {
		_, ok := roleCommitteeIdx[c.Committee]
		if !ok {
			missing = append(missing, c)
			continue
		}
	}

	log.Printf("found %d committees with P or O membership without Voter\n", len(missing))
	return missing
}
