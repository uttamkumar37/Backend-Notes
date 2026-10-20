package main

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"
)

const maxBodyBytes = 4096

type api struct {
	store *Store
}

func newHandler(store *Store) http.Handler {
	a := &api{store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/v1/investors", a.list)
	mux.HandleFunc("POST /api/v1/investors", a.create)
	mux.HandleFunc("GET /api/v1/investors/{id}", a.get)
	mux.HandleFunc("PUT /api/v1/investors/{id}", a.update)
	mux.HandleFunc("DELETE /api/v1/investors/{id}", a.delete)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// These response types are JSON-encodable. A write failure means the client
	// connection is gone; attempting a second HTTP response cannot repair it.
	_ = json.NewEncoder(w).Encode(value)
}

func fail(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func domainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		fail(w, http.StatusBadRequest, "name must contain 1–100 code points; country must contain two letters")
	case errors.Is(err, ErrNotFound):
		fail(w, http.StatusNotFound, "investor not found")
	default:
		fail(w, http.StatusInternalServerError, "internal error")
	}
}

func decode(w http.ResponseWriter, r *http.Request) (Input, bool) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		fail(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return Input{}, false
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var input *Input
	err = decoder.Decode(&input)
	if err == nil {
		var trailing any
		err = decoder.Decode(&trailing)
		if err == io.EOF && input != nil {
			return *input, true
		}
	}
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		fail(w, http.StatusRequestEntityTooLarge, "body exceeds 4096 bytes")
	} else {
		fail(w, http.StatusBadRequest, "body must be one JSON object with name and country")
	}
	return Input{}, false
}

func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		fail(w, http.StatusBadRequest, "id must be a positive integer")
		return 0, false
	}
	return id, true
}

func (a *api) create(w http.ResponseWriter, r *http.Request) {
	input, ok := decode(w, r)
	if !ok {
		return
	}
	item, err := a.store.Create(input)
	if err != nil {
		domainError(w, err)
		return
	}
	w.Header().Set("Location", "/api/v1/investors/"+strconv.FormatInt(item.ID, 10))
	writeJSON(w, http.StatusCreated, item)
}

func (a *api) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	item, err := a.store.Get(id)
	if err != nil {
		domainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *api) list(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fail(w, http.StatusBadRequest, "invalid query encoding")
		return
	}
	after, limit := int64(0), 20
	for key, values := range query {
		if (key != "after" && key != "limit") || len(values) != 1 {
			fail(w, http.StatusBadRequest, "use only one after and one limit parameter")
			return
		}
	}
	if query.Has("after") {
		value, err := strconv.ParseInt(query.Get("after"), 10, 64)
		if err != nil || value < 0 {
			fail(w, http.StatusBadRequest, "after must be a nonnegative integer")
			return
		}
		after = value
	}
	if query.Has("limit") {
		value, err := strconv.Atoi(query.Get("limit"))
		if err != nil || value < 1 || value > 100 {
			fail(w, http.StatusBadRequest, "limit must be between 1 and 100")
			return
		}
		limit = value
	}
	writeJSON(w, http.StatusOK, a.store.List(after, limit))
}

func (a *api) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	input, ok := decode(w, r)
	if !ok {
		return
	}
	item, err := a.store.Update(id, input)
	if err != nil {
		domainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *api) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := a.store.Delete(id); err != nil {
		domainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
