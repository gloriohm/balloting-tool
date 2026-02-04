package processing

import "ballot-tool/internal/models"

func createRoleComIdx(roles []models.Role) map[string]*models.Role {
	idx := make(map[string]*models.Role, len(roles))

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
