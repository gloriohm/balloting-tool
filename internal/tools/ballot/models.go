package ballot

import "time"

type Ballot struct {
	Committee string
	Reference string
	Closing   time.Time
	Title     string
	URL       string
}

type BallotMatched struct {
	Ballot Ballot
	Voter  Voter
}

type Voter struct {
	FirstName string
	LastName  string
	Email     string
}
