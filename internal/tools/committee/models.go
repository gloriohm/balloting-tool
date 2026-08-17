package committee

import "time"

type Engagement struct {
	Committee string
	FirstName string
	LastName  string
	Email     string
	From      time.Time
	To        *time.Time
}
