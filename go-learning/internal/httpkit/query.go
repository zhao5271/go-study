package httpkit

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var ErrInvalidQuery = errors.New("invalid query")

func ParsePageSize(r *http.Request) (page int, size int, err error) {
	page = 1
	size = 20

	q := r.URL.Query()
	if raw := strings.TrimSpace(q.Get("page")); raw != "" {
		n, convErr := strconv.Atoi(raw)
		if convErr != nil || n < 1 {
			return 0, 0, ErrInvalidQuery
		}
		page = n
	}
	if raw := strings.TrimSpace(q.Get("size")); raw != "" {
		n, convErr := strconv.Atoi(raw)
		if convErr != nil || n < 1 || n > 100 {
			return 0, 0, ErrInvalidQuery
		}
		size = n
	}
	return page, size, nil
}

