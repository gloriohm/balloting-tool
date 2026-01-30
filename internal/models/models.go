package models

import "time"

type Ballot struct {
	Source    string // "cen" or "iso"
	Committee string
	Reference string
	Closing   time.Time
	Title     string
	URL       string
}

type Role struct {
	Committee string
	FirstName string
	LastName  string
	Email     string
}

type BallotWithRole struct {
	Ballot Ballot
	Role   Role
}
