package filereader

import (
	"ballot-tool/internal/utils/normalization"
	"strings"
)

func normalizeHeaders(headers []string) []string {
	out := make([]string, len(headers))
	for i, s := range headers {
		out[i] = normalization.NormalizeString(s)
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
