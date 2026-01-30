package config

type Config struct {
	CentralizedVoters []string `json:"centralizedVoters"`
	OutputPath        string   `json:"outputPath"`
}
