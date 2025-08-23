package components

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Handler provides HTTP handlers for component operations.
type Handler struct {
	store *Store
}

// NewHandler creates a Handler with the given database.
func NewHandler(db *sql.DB) *Handler {
	return &Handler{store: NewStore(db)}
}

type createRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Components handles /components for listing and creation.
func (h *Handler) Components(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.create(w, r)
	case http.MethodGet:
		h.list(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ComponentByID handles /components/{id} for retrieval.
func (h *Handler) ComponentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/components/")
	c, err := h.store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, c, http.StatusOK)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Code == "" {
		http.Error(w, "name and code required", http.StatusBadRequest)
		return
	}
	if err := validateCode(req.Code); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := randomID()
	if err != nil {
		http.Error(w, "id generation failed", http.StatusInternalServerError)
		return
	}
	schema := extractPropsSchema(req.Code)
	comp := &Component{ID: id, Name: req.Name, Code: req.Code, PropsSchema: schema}
	if err := h.store.Create(r.Context(), comp); err != nil {
		errText := strings.ToLower(err.Error())
		if strings.Contains(errText, "duplicate") || strings.Contains(errText, "unique") {
			http.Error(w, "name already exists", http.StatusConflict)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, comp, http.StatusCreated)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	comps, err := h.store.List(r.Context())
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, comps, http.StatusOK)
}

func respondJSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
