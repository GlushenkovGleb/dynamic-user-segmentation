package model

type Group struct {
	ID       int    `gorm:"id"`
	ParentID *int   `gorm:"parent_id"`
	Name     string `gorm:"name"`
}

func (Group) TableName() string {
	return "groups"
}

type GroupWithStudent struct {
	GroupID   int `gorm:"group_id"`
	StudentID int `gorm:"student_id"`
}

func (GroupWithStudent) TableName() string {
	return "group_with_student"
}

type Student struct {
	ID      int    `gorm:"id"`
	Name    string `gorm:"name"`
	Email   string `gorm:"email"`
	GroupID int    `gorm:"group_id"`
}

func (Student) TableName() string {
	return "students"
}

type GroupFull struct {
	GroupID   int         `json:"group_id"`
	SubGroups []GroupFull `json:"subGroups"`
	Name      string      `json:"name"`
}
