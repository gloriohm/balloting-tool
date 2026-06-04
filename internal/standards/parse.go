package standards

import (
	"log"
	"strings"
)

func rowToStandard(rows []map[string]string) []Standard {
	out := make([]Standard, 0, len(rows))
	for _, row := range rows {
		ref := normalizeReference(row["reference"])
		out = append(out, Standard{
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
