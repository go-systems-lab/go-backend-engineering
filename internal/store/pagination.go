package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (q PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return q, nil
		}

		q.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return q, nil
		}

		q.Offset = o
	}

	sort := qs.Get("sort")
	if sort != "" {
		q.Sort = sort
	}

	tags := qs.Get("tags")
	if tags != "" {
		q.Tags = strings.Split(tags, ",")
	} else {
		q.Tags = []string{}
	}

	search := qs.Get("search")
	if search != "" {
		q.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		q.Since = parseTime(since)
	}

	until := qs.Get("until")
	if until != "" {
		q.Until = parseTime(until)
	}

	return q, nil
}

func parseTime(s string) string {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}
