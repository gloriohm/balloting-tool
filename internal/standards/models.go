package standards

type Standard struct {
	Reference string
	Language  string
	Title     string
	URN       string
}

type StandardExpanded struct {
	Reference string
	TitleNO   string
	TitleEN   string
	Committee string
	Year      int
}

type CommitteeData struct {
	ProjectManager  string
	CommitteeStatus string
}

type AktualitetData struct {
	StandardExpanded
	CommitteeData
}
