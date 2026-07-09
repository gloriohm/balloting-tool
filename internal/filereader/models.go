package filereader

import "time"

type Row map[string]string

type NationalEngagementsRow struct {
	Committee  Committee
	Commitment Commitment
	Person     Person
}

type Committee struct {
	Reference       string
	Title           string
	Domain          CommitteeDomain
	Level           CommitteeLevel
	Status          CommitteeStatus
	Established     time.Time
	MirrorCommittee bool
}

type Commitment struct {
	Role    string
	Status  CommitmentStatus
	Start   time.Time
	End     time.Time
	Company Company
}

type Person struct {
	FirstName string
	LastName  string
	Email     string
	Gender    Gender
	Company   Company
}

type Company struct {
	Name     string
	Category string
	OrgForm  string
}

type StandardDashboardRow struct {
	PubStatus string
	Language  string
	ImportID  string
	Title     string
	Reference string
	SDO       string
	SareptaID string
}

type StandardAdvancedRow struct {
}

type BallotRow struct {
	BallotType      string
	Committee       string
	BallotReference string
	OpeningDate     time.Time
	ClosingDate     time.Time
	BallotName      string
	URL             string
}

// strongly typed fields not necessary when switching to read from GD
// can be switched out for simpler checks

type CommitteeDomain string

const (
	CommitteeDomainNational      CommitteeDomain = "national"
	CommitteeDomainRegional      CommitteeDomain = "regional"
	CommitteeDomainInternational CommitteeDomain = "international"
)

type CommitteeLevel string

const (
	CommitteeLevelTC CommitteeLevel = "tc"
	CommitteeLevelSC CommitteeLevel = "sc"
	CommitteeLevelWG CommitteeLevel = "wg"
)

type CommitteeStatus string

const (
	CommitteeStatusActive     = "active"
	CommitteeStatusTerminated = "terminated"
	CommitteeStatusInactive   = "inactive"
	CommitteeStatusSuspended  = "suspended"
	CommitteeStatusInProgress = "in_progress"
)

type CommitmentStatus string

const (
	CommitmentStatusActive     = "active"
	CommitmentStatusTerminated = "terminated"
)

type Gender string

const (
	GenderMale    = "m"
	GenderFemale  = "f"
	GenderUnknown = "unknown/other"
)
