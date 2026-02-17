package security

import (
	"errors"
	"strconv"
	"strings"
)

func parseSessionIDFromPath(path string) (string, error) {
	trimmed := strings.TrimPrefix(path, "/api/v1/security/sessions/")
	if trimmed == "" {
		return "", errors.New("sessionId is required")
	}

	sessionID := strings.Split(trimmed, "/")[0]
	if sessionID == "" {
		return "", errors.New("sessionId is required")
	}

	return sessionID, nil
}

func parseEventIDFromPath(path string) (int64, error) {
	trimmed := strings.TrimPrefix(path, "/api/v1/security/events/")
	if trimmed == "" {
		return 0, errors.New("eventId is required")
	}

	idStr := strings.Split(trimmed, "/")[0]
	if idStr == "" {
		return 0, errors.New("eventId is required")
	}

	eventID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || eventID <= 0 {
		return 0, errors.New("invalid event ID in path")
	}

	return eventID, nil
}
