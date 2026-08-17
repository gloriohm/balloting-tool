package parsing

func ToSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))

	for _, v := range values {
		set[v] = struct{}{}
	}

	return set
}
