package sdimport

type Standard struct {
	Reference       string
	Title           []Title
	Abstract        string
	ICS             []string
	SusDevGoals     []string
	Edition         int
	Pages           int
	Developer       []string
	Owner           []string
	PublicationDate string
}

type Project struct {
	Reference        string            `json:"reference"`
	Title            []Title           `json:"title"`
	Abstract         []Abstract        `json:"abstract"`
	Classifications  []Classification  `json:"classifications"`
	ProjectRelations []ProjectRelation `json:"projectRelations"`
	Edition          string            `json:"edition"`
	Owner            CommitteeResponse `json:"owner"`
	Developer        CommitteeResponse `json:"developer"`
	PubLink          []URN             `json:"publications"`
}

type Publication struct {
	Reference       string        `json:"reference"`
	ProjectID       URN           `json:"project"`
	PublicationDate string        `json:"publicationDate"`
	ReleaseItems    []ReleaseItem `json:"releaseItems"`
}

type Abstract struct {
	Format   string `json:"format"`
	Content  string `json:"content"`
	Language string `json:"lang"`
}

type Title struct {
	Language string `json:"lang"`
	Value    string `json:"value"`
}

type Classification struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type CommitteeResponse struct {
	TargetType  string `json:"targetType"`
	DisplayName string `json:"displayName"`
	URN         string `json:"urn"`
}

type ReleaseItem struct {
	Type     string   `json:"type"`
	Format   string   `json:"format"`
	Pages    int      `json:"pages"`
	Language []string `json:"contentLanguage"`
}

type Response struct {
	Project     []Project     `json:"project"`
	Publication []Publication `json:"publication"`
}

type ProjectRelation struct {
	Type            string          `json:"type"`
	ExternalProject ExternalProject `json:"externalProject"`
	URN             string          `json:"urn"`
}

type ExternalProject struct {
	Originator  string `json:"originator"`
	DisplayName string `json:"displayName"`
	ProjectID   string `json:"projectId"`
}

type URN struct {
	URN string `json:"urn"`
}

type Committee struct {
	NSB       string
	Reference string
}

var sustainableDevelopmentGoals = map[int]string{
	1:  "1: No Poverty",
	2:  "2: Zero Hunger",
	3:  "3: Good Health and Well-Being",
	4:  "4: Quality Education",
	5:  "5: Gender Equality",
	6:  "6: Clean Water and Sanitation",
	7:  "7: Affordable and Clean Energy",
	8:  "8: Decent Work and Economic Growth",
	9:  "9: Industry, Innovation and Infrastructure",
	10: "10: Reduced Inequalities",
	11: "11: Sustainable Cities and Communities",
	12: "12: Responsible Consumption and Production",
	13: "13: Climate Action",
	14: "14: Life Below Water",
	15: "15: Life on Land",
	16: "16: Peace, Justice and Strong Institutions",
	17: "17: Partnerships for the Goals",
}
