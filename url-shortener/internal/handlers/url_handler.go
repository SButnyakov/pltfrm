package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"url-shortener/internal/models"
)

type urlService interface {
	Create(address string) (string, error)
	GetByURL(url string) (*models.URL, error)
}

type urlHandler struct {
	service urlService
}

func NewUrlHandler(service urlService) *urlHandler {
	return &urlHandler{
		service: service,
	}
}

func (h *urlHandler) CreateURL(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Address string `json:"address"`
	}
	var req request

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		slog.Error("failed to decode JSON from body", slog.Any("error", err), slog.Any("body", r.Body))
		render.Status(r, http.StatusBadRequest)
		return
	}

	render.Status(r, http.StatusCreated)
	url, err := h.service.Create(req.Address)
	if err != nil {
		slog.Error("failed to create URL", slog.Any("error", err))
		render.Status(r, http.StatusInternalServerError)
		return
	}
	url = fmt.Sprintf("http://localhost:8080/%s", url)

	type response struct {
		URL string `json:"url"`
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response{url})
}

func (h *urlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")

	model, err := h.service.GetByURL(url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			render.Status(r, http.StatusNotFound)
			return
		}
		render.Status(r, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, model.Address, http.StatusTemporaryRedirect)
}
