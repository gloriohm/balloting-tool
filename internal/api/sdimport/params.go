package sdimport

import "fmt"

type Parameters struct {
	Vendor              string
	PageSize            int
	FromDate            string
	ToDate              string
	LastChangeTimestamp string
	Originator          string
}

func NewParameters(from, to string) Parameters {
	return Parameters{
		Vendor:              "sarepta",
		PageSize:            50,
		FromDate:            from,
		ToDate:              to,
		LastChangeTimestamp: "2023-01-01T00:00:00",
		Originator:          "SN",
	}
}

func (p *Parameters) buildRequestString(pubType string, page int) string {
	return fmt.Sprintf("/%s/%s/%s/%d/%d?publicationDateFrom=%s&publicationDateTo=%s&originator=%s", pubType, p.Vendor, p.LastChangeTimestamp, page, p.PageSize, p.FromDate, p.ToDate, p.Originator)
}
