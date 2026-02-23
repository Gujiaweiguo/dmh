package handlerutil

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ParsePathID(r *http.Request, w http.ResponseWriter, prefix string, errorMsg string) (int64, bool) {
	path := strings.TrimPrefix(r.URL.Path, prefix)
	idStr := strings.Split(path, "/")[0]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, errors.New(errorMsg))
		return 0, false
	}
	return id, true
}

func ParseQueryInt64(r *http.Request, key string, defaultValue int64) int64 {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil || val <= 0 {
		return defaultValue
	}
	return val
}

func ParseQueryString(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
