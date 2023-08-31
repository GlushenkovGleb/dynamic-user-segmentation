package model

import "time"

type MemberEvent struct {
	UserId     string    `db:"user_id"`
	Slug       string    `db:"slug"`
	Status     string    `db:"status"`
	HappenedAt time.Time `db:"happened_at"`
}

type SegmentToAdd struct {
	Slug         string `json:"slug"`
	TttInSeconds int    `json:"ttl_in_seconds,omitempty" validate:"gte=0"`
}

type Segment struct {
	Slug       string  `json:"slug" validate:"required"`
	Percentage float64 `json:"percentage,omitempty"`
	IsAuto     bool    `json:"is_auto,omitempty"`
}

type UserSegments struct {
	Slugs []string `json:"slugs"`
}
