package ballot

import (
	"ballot-tool/internal/filereader"
)

func rowToVoter(in filereader.Person) Voter {
	return Voter{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
	}
}

func rowToBallot(in filereader.BallotRow) Ballot {
	return Ballot{
		Committee: in.Committee,
		Reference: in.BallotReference,
		Closing:   in.ClosingDate,
		Title:     in.BallotName,
		URL:       in.URL,
	}
}
