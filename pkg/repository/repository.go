package repository

import (
	"context"
	"errors"

	"dynamic-user-segmentation/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const path = "storage.db"

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	query := `
	

	Create Table IF NOT EXISTS groups (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name varchar(255) NOT NULL,
	    parent_id INTEGER REFERENCES groups(id)
	);
	
	CREATE TABLE IF NOT EXISTS students (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name varchar(255) NOT NULL,
	    email varchar(255) NOT NULL,
	    group_id INTEGER REFERENCES groups(id)
	);
	`
	err = db.Exec(query).Error
	if err != nil {
		return nil, err
	}

	return db, nil
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateGroup(ctx context.Context, group *model.Group) error {
	r.db.Create(group)
	return r.db.Error
}

func (r *Repository) GetGroups(ctx context.Context) ([]*model.GroupFull, error) {
	parentGroups := make([]*model.Group, 0)
	err := r.db.Where("parent_id is NULL").Find(&parentGroups).Error
	if err != nil {
		return nil, err
	}

	if len(parentGroups) == 0 {
		return nil, nil
	}

	fullGroups := make([]*model.GroupFull, 0)
	for _, parentGroup := range parentGroups {
		children, err := r.getChildren(ctx, parentGroup.ID)
		if err != nil {
			return nil, err
		}

		fullGroups = append(fullGroups, &model.GroupFull{
			GroupID:   parentGroup.ID,
			Name:      parentGroup.Name,
			SubGroups: children,
		})
	}

	return fullGroups, nil
}

func (r *Repository) getChildren(ctx context.Context, parentID int) ([]model.GroupFull, error) {
	groups := make([]model.Group, 0)
	err := r.db.Find(&groups, "parent_id = ?", parentID).Error
	if err != nil {
		return nil, err
	}

	fullGroups := make([]model.GroupFull, 0)
	for _, group := range groups {
		children, err := r.getChildren(ctx, group.ID)
		if err != nil {
			return nil, err
		}

		fullGroups = append(fullGroups, model.GroupFull{
			GroupID:   group.ID,
			Name:      group.Name,
			SubGroups: children,
		})
	}

	return fullGroups, nil
}

func (r *Repository) GetGroup(ctx context.Context, groupID int) (*model.Group, error) {
	group := &model.Group{}
	err := r.db.Where("id = ?", groupID).First(group).Error
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (r *Repository) UpdateGroup(ctx context.Context, group *model.Group) error {
	return r.db.Updates(group).Error
}

func (r *Repository) GetGroupsByName(ctx context.Context, name string) ([]model.Group, error) {
	var groups []model.Group
	err := r.db.Where("name = ?", name).Find(&groups).Error
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *Repository) DeleteGroup(ctx context.Context, groupID int) error {
	err := r.db.Where("parent_id = ?", groupID).First(&model.Group{}).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Where("parent_id = ?", groupID).Delete(&model.Group{}).Error
	}

	return errors.New("cannot delete group")
}

func (r *Repository) CreateStudent(ctx context.Context, student *model.Student) error {
	return r.db.Create(student).Error
}

func (r *Repository) GetStudents(ctx context.Context) ([]model.Student, error) {
	var students []model.Student
	err := r.db.Find(&students).Error
	if err != nil {
		return nil, err
	}

	return students, nil
}

func (r *Repository) GetStudent(ctx context.Context, id int) (model.Student, error) {
	var student model.Student
	err := r.db.Where("id = ?", id).First(&student).Error
	if err != nil {
		return model.Student{}, err
	}

	return student, nil
}

func (r *Repository) UpdateStudent(ctx context.Context, student *model.Student) error {
	return r.db.Updates(student).Error
}

func (r *Repository) GetStudentsByGroupName(ctx context.Context, groupName string) ([]model.Student, error) {
	var students []model.Student
	err := r.db.Select("students.name, students.id, students.group_id").Joins("join groups g on g.id=students.group_id").Where("g.name = ?", groupName).Find(&students).Error
	if err != nil {
		return nil, err
	}

	return students, nil
}

func (r *Repository) DeleteStudent(ctx context.Context, id int) error {
	return r.db.Where("id = ?", id).Delete(&model.Student{}).Error
}
