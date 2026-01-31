package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/aftaab/trelay/internal/api/response"
	"github.com/aftaab/trelay/internal/core/domain"
	"github.com/aftaab/trelay/internal/core/folder"
)

type FolderHandler struct {
	service *folder.Service
}

func NewFolderHandler(service *folder.Service) *FolderHandler {
	return &FolderHandler{service: service}
}

func (h *FolderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	createdFolder, err := h.service.Create(r.Context(), req)
	if err != nil {
		if ve, ok := err.(domain.ValidationError); ok {
			response.ValidationError(w, ve.Field, ve.Message)
			return
		}
		if err == domain.ErrParentFolderNotFound {
			response.ValidationError(w, "parent_id", "parent folder not found")
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, createdFolder)
}

func (h *FolderHandler) List(w http.ResponseWriter, r *http.Request) {
	folders, err := h.service.List(r.Context())
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, folders)
}

func (h *FolderHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid folder ID")
		return
	}

	f, err := h.service.Get(r.Context(), id)
	if err != nil {
		if err == domain.ErrFolderNotFound {
			response.NotFound(w, "folder not found")
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, f)
}

func (h *FolderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid folder ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == domain.ErrFolderNotFound {
			response.NotFound(w, "folder not found")
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"deleted": true})
}
