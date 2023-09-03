package repository

import (
	"dynamic-user-segmentation/pkg/model"
	"github.com/jmoiron/sqlx"
)

const deletedStatus = "DELETED"

const insertSegmentQuery = "INSERT INTO segments(slug, percentage, is_auto) VALUES ($1, $2, $3) RETURNING id"

const selectDeleteEventsBySlugQuery = `
	SELECT user_id, slug, 
	CASE 
		WHEN expires_at is null THEN now()
		WHEN now() > expires_at THEN expires_at
		ELSE now()
	END AS happened_at
	FROM segments
	INNER JOIN test_members ON segments.id = test_members.segment_id
	WHERE slug = $1
`

const insertDeleteEventsQuery = `
	INSERT INTO test_members_history(user_id, slug, status, happened_at)
	VALUES (:user_id, :slug, :status, :happened_at)
`

const deleteAllMembersBySlugQuery = `
	DELETE  FROM test_members
	WHERE segment_id IN (SELECT id FROM segments WHERE slug = $1)
`

const deleteSegmentBySlugQuery = `
	DELETE FROM segments
	WHERE slug = $1
	RETURNING id
`

type SegmentPostgres struct {
	db *sqlx.DB
}

func NewSegmentPostgres(db *sqlx.DB) *SegmentPostgres {
	return &SegmentPostgres{db: db}
}

func (s *SegmentPostgres) CreateSegment(segment model.Segment) (int, error) {
	var id int
	row := s.db.QueryRow(insertSegmentQuery, segment.Slug, segment.Percentage, segment.IsAuto)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *SegmentPostgres) DeleteSegment(slug string) (int, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, nil
	}
	defer tx.Rollback()

	// Берем записи об участии в сегменте
	var deleteEvents []model.MemberEvent
	err = tx.Select(&deleteEvents, selectDeleteEventsBySlugQuery, slug)
	if err != nil {
		return 0, err
	}

	// Сохраняем записи об участии в историю
	for i := range deleteEvents {
		deleteEvents[i].Status = deletedStatus
	}
	if len(deleteEvents) > 0 {
		_, err = tx.NamedExec(insertDeleteEventsQuery, deleteEvents)
		if err != nil {
			return 0, err
		}
	}

	// Удаляем все записи с участием сегмента
	_, err = tx.Exec(deleteAllMembersBySlugQuery, slug)
	if err != nil {
		return 0, err
	}

	// Удаляем сегмент
	res, err := tx.Exec(deleteSegmentBySlugQuery, slug)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
