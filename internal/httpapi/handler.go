package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Bazu-Bazu/Http-Server-On-Go/internal/user"
)

type UserService interface {
	Register(input user.RegisterInput) (user.User, error)
	Get(id int64) (user.User, error)
	Delete(id int64) error
}

type Handler struct {
	users UserService
}

func NewHandler(users UserService) *Handler {
	return &Handler{users: users}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("POST /api/users", h.register)
	mux.HandleFunc("GET /api/users/{id}", h.getByID)
	mux.HandleFunc("DELETE /api/users/{id}", h.deleteByID)
	return mux
}

type registerRequest struct {
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type userResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
}

func toUserResponse(u user.User) userResponse {
	return userResponse{
		ID:       u.ID,
		Nickname: u.Nickname,
		Name:     u.Name,
	}
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_json", "request body must be one valid JSON object with known fields")
		return
	}

	u, err := h.users.Register(user.RegisterInput{
		Nickname: req.Nickname,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toUserResponse(u))
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	u, err := h.users.Get(id)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(u))
}

func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.users.Delete(id); err != nil {
		h.writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 0 {
		writeAPIError(w, http.StatusBadRequest, "invalid id", "id must be a positive integer")
		return 0, false
	}

	return id, true
}

func (h *Handler) writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, user.ErrInvalidInput):
		writeAPIError(w, http.StatusBadRequest, "invalid_input", err.Error())
	case errors.Is(err, user.ErrNotFound):
		writeAPIError(w, http.StatusNotFound, "user_not_found", "user was not found")
	case errors.Is(err, user.ErrNicknameTaken):
		writeAPIError(w, http.StatusConflict, "nickname_taken", "this nickname is already in use")
	default:
		log.Printf("internal service error: %v", err)
		writeAPIError(w, http.StatusInternalServerError, "internal_error", "internal server error")
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return err
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("request contains more than one JSON value")
		}
		return err
	}
	return nil
}
