package config

type Config struct {
	CentralizedVoters []string `json:"centralizedVoters"`
	OutputPath        string   `json:"outputPath"`
	InputPath         string   `json:"inputPathSuffix"`
	Files             Files    `json:"files"`
}

type Files struct {
	Ballot1  string `json:"ballot1"`
	Ballot2  string `json:"ballot2"`
	Voters   string `json:"voterRoles"`
	OrgRoles string `json:"orgRoles"`
}
