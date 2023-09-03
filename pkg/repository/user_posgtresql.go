package repository

import (
	"dynamic-user-segmentation/pkg/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

const selectActiveUserSegmentsQuery = `
	SELECT slug
	FROM test_members
	INNER JOIN segments ON test_members.segment_id = segments.id 
	WHERE user_id = $1
	AND (expires_at is null OR expires_at > now())
`

const insertUserIfNotExistsQuery = `
	INSERT INTO users(id)
	SELECT CAST($1 AS varchar)
	WHERE NOT EXISTS (
		SELECT id FROM users WHERE id = $1
	)
`

const saveNewAddedMembersQuery = `
	INSERT INTO test_members_history(user_id, slug, status, happened_at)
	SELECT :userId, :slug, 'ADDED', now()
	WHERE NOT EXISTS (
		SELECT tm.id
		FROM test_members tm
		INNER JOIN segments ON segments.id = tm.segment_id
		WHERE user_id = :userId AND slug = :slug
	)
`

const upsertUserSegmentsWithTTLQuery = `
	INSERT INTO test_members(user_id, segment_id, expires_at)
	SELECT :userId, seg.id,
	CASE
		WHEN :ttlInSeconds = 0 THEN NULL
		ELSE now() + :ttlInSeconds * interval '1 second'
	END AS expires_at
	FROM segments seg
	WHERE slug = :slug
	ON CONFLICT ON CONSTRAINT uc_member
	DO NOTHING 
`

const saveDeletedMembersQuery = `
	INSERT INTO test_members_history(user_id, slug, status, happened_at)
	SELECT :userId, :slug, 'DELETED',
	CASE
		WHEN expires_at is null THEN now()
		WHEN now() > expires_at THEN expires_at
		ELSE now()
	END AS happened_at
	FROM test_members
	INNER JOIN segments ON test_members.segment_id = segments.id
	WHERE user_id = :userId AND slug = :slug
`

const deleteMembersRecordQuery = `
	DELETE FROM test_members
	WHERE user_id = :userId AND
	      segment_id IN (SELECT id FROM segments WHERE slug = :slug) 
`

const saveExpiredMembersHistoryQuery = `
	INSERT INTO test_members_history(user_id, slug, status, happened_at)
	SELECT user_id, slug, 'DELETED', tm.expires_at
	FROM test_members tm
	INNER JOIN segments sg ON  sg.id = tm.segment_id
	WHERE expires_at > now()
`

const deleteExpiredMembersHistoryQuery = `
	DELETE FROM test_members
	WHERE expires_at > now()
`

const getMembersMonthHistoryQuery = `
	SELECT user_id, slug, status, happened_at
	FROM test_members_history
	WHERE extract(year from happened_at) = $1 AND extract(month from happened_at) = $2
`

func (u *UserPostgres) GetSegments(userId string) (model.UserSegments, error) {
	segments := model.UserSegments{Slugs: []string{}}
	err := u.db.Select(&segments.Slugs, selectActiveUserSegmentsQuery, userId)
	if err != nil {
		return segments, err
	}

	return segments, nil
}

func (u *UserPostgres) GetSegmentsHistory(year, month int) ([]model.MemberEvent, error) {
	tx, err := u.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Сохраняем истекших участников
	_, err = tx.Exec(saveExpiredMembersHistoryQuery)
	if err != nil {
		return nil, err
	}

	// Удаляем истекших участников
	_, err = tx.Exec(deleteExpiredMembersHistoryQuery)
	if err != nil {
		return nil, err
	}

	// Получаем историю за определенные период
	events := make([]model.MemberEvent, 0)
	err = tx.Select(&events, getMembersMonthHistoryQuery, year, month)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (u *UserPostgres) UpdateSegments(userId string, segmentsToAdd []model.SegmentToAdd, slugsToDelete []string) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Убеждаемся, что пользователь есть в базе
	_, err = tx.Exec(insertUserIfNotExistsQuery, userId)
	if err != nil {
		return err
	}
	fmt.Println(1)

	if len(slugsToDelete) > 0 {
		// Сохраняем в историю события об удалении сегментов
		deletedParams := prepareDeletedParams(userId, slugsToDelete)
		_, err = tx.NamedExec(saveDeletedMembersQuery, deletedParams)
		if err != nil {

		}
		fmt.Println(2)

		// Удаляем пользовательские сегменты
		_, err = tx.NamedExec(deleteMembersRecordQuery, deletedParams)
		if err != nil {
			return err
		}
		fmt.Println(3)
	}

	newUserSegments := prepareUserSegments(userId, segmentsToAdd)
	if len(newUserSegments) > 0 {
		// Сохраняем в историю добавление новых сегментов
		_, err = tx.NamedExec(saveNewAddedMembersQuery, newUserSegments)
		if err != nil {
			return err
		}
		fmt.Println(2)

		// Добавляем пользователю сегменты
		_, err = tx.NamedExec(upsertUserSegmentsWithTTLQuery, newUserSegments)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func prepareUserSegments(userId string, segmentsToAdd []model.SegmentToAdd) []map[string]interface{} {
	newUserSegments := make([]map[string]interface{}, 0, len(segmentsToAdd))
	for _, segment := range segmentsToAdd {
		upsertParam := map[string]interface{}{
			"slug":         segment.Slug,
			"ttlInSeconds": segment.TttInSeconds,
			"userId":       userId,
		}
		newUserSegments = append(newUserSegments, upsertParam)
	}
	return newUserSegments
}

func prepareDeletedParams(userId string, slugsToDelete []string) []map[string]interface{} {
	deletedParams := make([]map[string]interface{}, 0, len(slugsToDelete))
	for _, slug := range slugsToDelete {
		deleteParam := map[string]interface{}{
			"slug":   slug,
			"userId": userId,
		}
		deletedParams = append(deletedParams, deleteParam)
	}
	return deletedParams
}
