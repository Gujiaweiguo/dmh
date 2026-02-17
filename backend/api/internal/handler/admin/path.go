package admin

import (
	"errors"
	"strconv"
	"strings"
)

func parseUserIDFromPath(path string) (int64, error) {
	trimmed := strings.TrimPrefix(path, "/api/v1/admin/users/")
	if trimmed == "" {
		return 0, errors.New("userId is required")
	}

	idStr := strings.Split(trimmed, "/")[0]
	if idStr == "" {
		return 0, errors.New("userId is required")
	}

	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || userId <= 0 {
		return 0, errors.New("invalid user ID in path")
	}

	return userId, nil
}
