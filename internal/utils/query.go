package utils

import "net/http"

func DefaultQuery(r *http.Request, p string, d string) string {
	q := r.URL.Query()

	pv := q.Get(p)
	if pv == "" {
		pv = d
	}

	return pv
}
