package filereader

import (
	"ballot-tool/internal/utils/normalization"
	"strings"
)

func normalizeHeaders(headers []string) []string {
	out := make([]string, len(headers))
	for i, header := range headers {
		header = normalization.NormalizeString(header)
		switch header {
		case "opening_date":
			return "start"
		}
	}

	return out
}

func splitBySeparator(s, sep string) []string {
	parts := strings.Split(s, sep)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	return parts
}
