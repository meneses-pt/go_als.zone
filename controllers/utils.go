package controllers

import "net/http"

// HandleLimitOffsetParams gets the offset and limit parameters from the request string
func HandleLimitOffsetParams(r *http.Request) (string, string) {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "50"
	}
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}
	return limit, offset
}
