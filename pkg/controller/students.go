package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"dynamic-user-segmentation/pkg/model"
)

type CreateStudentRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	GroupID int    `json:"group_id"`
}

func (c *Controller) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var req CreateStudentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.repo.CreateStudent(r.Context(), &model.Student{
		ID:      req.GroupID,
		Name:    req.Name,
		Email:   req.Email,
		GroupID: req.GroupID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type GetStudentResp struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	GroupID int    `json:"group_id"`
}

func (c *Controller) GetStudents(w http.ResponseWriter, r *http.Request) {
	groupNameParam := r.URL.Query().Get("name")
	if groupNameParam == "" {
		students, err := c.repo.GetStudentsByGroupName(r.Context(), groupNameParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var studentsResp []GetStudentResp
		for _, student := range students {
			studentsResp = append(studentsResp, GetStudentResp{
				ID:      student.ID,
				Name:    student.Name,
				GroupID: student.GroupID,
			})
		}
		err = json.NewEncoder(w).Encode(studentsResp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		return
	}

	students, err := c.repo.GetStudents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var studentsResp []GetStudentResp
	for _, student := range students {
		studentsResp = append(studentsResp, GetStudentResp{
			ID:      student.ID,
			Name:    student.Name,
			GroupID: student.GroupID,
		})
	}
	err = json.NewEncoder(w).Encode(studentsResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) GetStudent(w http.ResponseWriter, r *http.Request) {
	studentIDParam := r.URL.Query().Get("student_id")
	studentID, err := strconv.Atoi(studentIDParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	student, err := c.repo.GetStudent(r.Context(), studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	studentResp := GetStudentResp{
		ID:      student.ID,
		Name:    student.Name,
		GroupID: student.GroupID,
	}
	err = json.NewEncoder(w).Encode(studentResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type UpdateStudentRequest struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	GroupID int    `json:"group_id"`
}

func (c *Controller) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	studentIDParam := r.URL.Query().Get("student_id")
	studentID, err := strconv.Atoi(studentIDParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req UpdateStudentRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.repo.UpdateStudent(r.Context(), &model.Student{
		ID:      studentID,
		Name:    req.Name,
		GroupID: req.GroupID,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	studentIDParam := r.URL.Query().Get("student_id")
	studentID, err := strconv.Atoi(studentIDParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.repo.DeleteStudent(r.Context(), studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
