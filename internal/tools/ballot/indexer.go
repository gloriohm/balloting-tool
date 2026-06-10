package ballot

func createRoleComIdx(roles []Role) map[string]*Role {
	idx := make(map[string]*Role, len(roles))

	for i := range roles {
		c := roles[i].Committee
		if c == "" {
			continue
		}
		if _, exists := idx[c]; !exists {
			idx[c] = &roles[i]
		}
	}

	return idx
}
