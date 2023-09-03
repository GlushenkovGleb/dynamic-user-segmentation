package repository

import (
	"dynamic-user-segmentation/pkg/model"
	"errors"
	"github.com/jmoiron/sqlx"
	"io"
)

var (
	ErrConnectError = "connection to db failed"
)

var (
	ErrSegmentNotFound = "segment doesn't exist"
)

type User interface {
	GetSegments(userId string) (model.UserSegments, error)
	GetSegmentsHistory(year, month int) ([]model.MemberEvent, error)
	UpdateSegments(userId string, segmentsToAdd []model.SegmentToAdd, segmentSlugsToDelete []string) error
}

type History interface {
	SaveEvents([]model.MemberEvent) (string, error)
	UploadEvents(fileName string, dest io.Writer) (int, error)
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
	History
}

func NewRepository(db *sqlx.DB, fileStoragePath string) *Repository {
	return &Repository{
		User:    NewUserPostgres(db),
		Segment: NewSegmentPostgres(db),
		History: NewHistoryCSV(fileStoragePath),
	}
}
