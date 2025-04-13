package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"dynamic-user-segmentation/pkg/model"
	"github.com/go-chi/chi/v5"
)

type CreateGroupRequest struct {
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
}

func (c *Controller) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req CreateGroupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group := model.Group{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	err = c.repo.CreateGroup(r.Context(), &group)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(group)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type GroupByName struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func (c *Controller) GetGroups(w http.ResponseWriter, r *http.Request) {
	nameQueryParam := r.URL.Query().Get("name")
	if nameQueryParam == "" {
		fullGroups, err := c.repo.GetGroups(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(fullGroups)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	} else {
		groups, err := c.repo.GetGroupsByName(r.Context(), nameQueryParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var groupsByName []GroupByName
		for _, group := range groups {
			groupsByName = append(groupsByName, GroupByName{
				Name: group.Name,
				ID:   group.ID,
			})
		}
		err = json.NewEncoder(w).Encode(groupsByName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type GetGroupResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Controller) GetGroup(w http.ResponseWriter, r *http.Request) {
	groupIDParam := chi.URLParam(r, "group_id")
	groupID, err := strconv.Atoi(groupIDParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	group, err := c.repo.GetGroup(r.Context(), groupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := GetGroupResponse{
		ID:   group.ID,
		Name: group.Name,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

type UpdateGroupRequest struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
}

func (c *Controller) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	var req UpdateGroupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	groupIDParam := chi.URLParam(r, "group_id")
	groupID, err := strconv.Atoi(groupIDParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.repo.UpdateGroup(r.Context(), &model.Group{
		ID:       groupID,
		Name:     req.Name,
		ParentID: req.ParentID,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupIDParam := chi.URLParam(r, "group_id")
	groupID, err := strconv.Atoi(groupIDParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.repo.DeleteGroup(r.Context(), groupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
