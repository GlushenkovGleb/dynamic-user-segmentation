package repository

import (
	"dynamic-user-segmentation/pkg/model"
	"errors"
	"github.com/jmoiron/sqlx"
)

var (
	ErrConnectError = "connection to db failed"
)

var (
	ErrSegmentNotFound = "segment doesn't exist"
)

type User interface {
	GetSegments(userId string) (model.UserSegments, error)
	GetSegmentsHistory() ([]model.MemberEvent, error)
	UpdateSegments(userId string, segmentsToAdd []model.SegmentToAdd, segmentSlugsToDelete []string) error
}

var (
	ErrSegmentExists     = errors.New("segment already exists in db")
	ErrSegmentValidation = errors.New("segment validation failed in db")
)

type Segment interface {
	CreateSegment(segment model.Segment) (int, error)
	DeleteSegment(slug string) (int, error)
}

type Repository struct {
	User
	Segment
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:    NewUserPostgres(db),
		Segment: NewSegmentPostgres(db),
	}
}
