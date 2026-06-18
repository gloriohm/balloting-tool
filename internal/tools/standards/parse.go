package standards

import (
	"ballot-tool/internal/api/sdimport"
	"log"
	"strings"
)

func rowToStandard(rows []map[string]string) []StandardCore {
	out := make([]StandardCore, 0, len(rows))
	for _, row := range rows {
		ref := normalizeReference(row["reference"])
		out = append(out, StandardCore{
			Reference: ref,
			Language:  row["lang"],
			Title:     row["title"],
			URN:       row["id"],
		})
	}

	log.Printf("parsed standards, length %d\n", len(out))
	return out
}

func normalizeReference(ref string) string {
	ref = strings.ReplaceAll(ref, "ISO/", "ISO_")
	ref = strings.ReplaceAll(ref, "CEN/", "CEN_")
	ref = strings.ReplaceAll(ref, "SN/", "SN_")
	return ref
}

func projToAktualitetStandard(proj sdimport.Project) AktualitetStandard {
	var out AktualitetStandard
	out.Reference = proj.Reference
	titles := proj.ParseTitles()
	for _, t := range titles {
		switch t.Language {
		case "no":
			out.TitleNO = t.Value
		case "en":
			out.TitleEN = t.Value
		}
	}
	out.Committee = proj.Owner.DisplayName
	year, _ := getYearFromReference(proj.Reference)
	out.Year = year

	return out
}

func pdfTarget() Target {
	return Target{itemType: sdimport.ReleaseItemTypeStandard, itemFormat: sdimport.ReleaseItemFormatPDF}
}

func xmlTarget() Target {
	return Target{itemType: sdimport.ReleaseItemTypeStandard, itemFormat: sdimport.ReleaseItemFormatXML}
}

func wordTarget() Target {
	return Target{itemType: sdimport.ReleaseItemTypeSource, itemFormat: sdimport.ReleaseItemFormatWord}
}

func anyTarget() Target {
	return Target{itemType: sdimport.ReleaseItemTypeOther, itemFormat: sdimport.ReleaseItemFormatAny}
}
