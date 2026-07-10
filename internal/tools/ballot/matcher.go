package ballot

import (
	"ballot-tool/internal/filereader"
	"ballot-tool/internal/utils/normalization"
	"log"
)

func matchBallotVoter(ballots []filereader.BallotRow, voters []filereader.NationalEngagementsRow) ([]BallotMatched, []filereader.BallotRow) {
	votersIdx := indexVoters(voters)
	matches := make([]BallotMatched, 0, len(ballots))
	var missing []filereader.BallotRow

	for _, b := range ballots {
		c := normalization.NormalizeString(b.Committee)
		match, ok := votersIdx[c]
		if !ok {
			missing = append(missing, b)
			continue
		}

		matches = append(matches, BallotMatched{
			Ballot: rowToBallot(b),
			Voter:  match,
		})
	}

	log.Printf("matched %d ballots \n", len(matches))
	log.Printf("found %d ballots without voter \n", len(missing))

	return matches, missing
}

func indexVoters(roles []filereader.NationalEngagementsRow) map[string]Voter {
	idx := make(map[string]Voter, len(roles))

	for i := range roles {
		c := normalization.NormalizeString(roles[i].Committee.Reference)
		if c == "" {
			continue
		}

		// extra check just in case something goes wrong with the intial filtering
		v := normalization.NormalizeString(roles[i].Commitment.Role)
		if v != "voter" {
			continue
		}

		if _, exists := idx[c]; !exists {
			idx[c] = rowToVoter(roles[i].Person)
		}
	}

	return idx
}
