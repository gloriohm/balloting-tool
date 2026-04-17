package ballot

import (
	"log"
)

func JoinBallotRole(roles []Role, ballots []Ballot) ([]BallotWithRole, []Ballot) {
	roleCommitteeIdx := createRoleComIdx(roles)

	matches := make([]BallotWithRole, 0, len(ballots))
	missing := make([]Ballot, 0, len(ballots))
	for _, b := range ballots {
		match, ok := roleCommitteeIdx[b.Committee]
		if !ok {
			missing = append(missing, b)
			continue
		}

		matches = append(matches, BallotWithRole{
			Ballot: b,
			Role:   *match,
		})
	}

	log.Printf("matched %d ballots \n", len(matches))
	log.Printf("found %d ballots without voter \n", len(missing))

	return matches, missing
}

func JoinCommitteeRole(roles []Role, coms []Committee) []Committee {
	roleCommitteeIdx := createRoleComIdx(roles)

	missing := make([]Committee, 0, len(coms))

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
