package committee

import (
	"ballot-tool/internal/filereader"
	"encoding/csv"
	"io"
	"os"
)

func loadAllCompanies(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)

	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	idx := make(map[string]int)
	for i, h := range header {
		idx[h] = i
	}
	result := make(map[string]string)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if record[idx["commitment_status"]] != "active" {
			continue
		}

		result[record[idx["orgnummer"]]] = record[idx["company_name"]]
	}

	return result, nil
}

func parseEngagement(in filereader.NationalEngagementsRow) Engagement {
	out := Engagement{
		Committee: in.Committee.Reference,
		FirstName: in.Person.FirstName,
		LastName:  in.Person.LastName,
		Email:     in.Person.Email,
		From:      in.Commitment.Start,
	}

	if !in.Commitment.End.IsZero() {
		out.To = &in.Commitment.End
	}

	return out
}

func parseEngagements(in []filereader.NationalEngagementsRow) []Engagement {
	var out []Engagement
	for _, e := range in {
		out = append(out, parseEngagement(e))
	}

	return out
}
