package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func request(handler http.Handler, method, path, body, contentType string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	if contentType != "" {
		r.Header.Set("Content-Type", contentType)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
}

func TestCRUDOverHTTP(t *testing.T) {
	// A real HTTP connection also exercises routing and response serialization.
	server := httptest.NewServer(newHandler(newStore()))
	defer server.Close()
	response, err := server.Client().Post(server.URL+"/api/v1/investors", "application/json", strings.NewReader(`{"name":" Asha ","country":"in"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated || response.Header.Get("Location") != "/api/v1/investors/1" {
		t.Fatalf("create: %s, location %q", response.Status, response.Header.Get("Location"))
	}
	var item Investor
	if err := json.NewDecoder(response.Body).Decode(&item); err != nil || item.Name != "Asha" || item.Country != "IN" {
		t.Fatalf("create body = %+v, %v", item, err)
	}
	got, err := server.Client().Get(server.URL + response.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get: %s", got.Status)
	}
}

func TestUpdateDeleteAndPagination(t *testing.T) {
	handler := newHandler(newStore())
	for _, name := range []string{"Asha", "Mira", "Ravi"} {
		w := request(handler, "POST", "/api/v1/investors", `{"name":"`+name+`","country":"IN"}`, "application/json")
		if w.Code != http.StatusCreated {
			t.Fatalf("seed: %d %s", w.Code, w.Body)
		}
	}
	w := request(handler, "GET", "/api/v1/investors?after=1&limit=1", "", "")
	var page []Investor
	if err := json.Unmarshal(w.Body.Bytes(), &page); err != nil || w.Code != http.StatusOK || len(page) != 1 || page[0].ID != 2 {
		t.Fatalf("page = %s, status %d, error %v", w.Body, w.Code, err)
	}
	w = request(handler, "PUT", "/api/v1/investors/2", `{"name":"Mira Updated","country":"GB"}`, "application/json")
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "Mira Updated") {
		t.Fatalf("update = %d %s", w.Code, w.Body)
	}
	w = request(handler, "DELETE", "/api/v1/investors/2", "", "")
	if w.Code != http.StatusNoContent || w.Body.Len() != 0 {
		t.Fatalf("delete = %d %s", w.Code, w.Body)
	}
	w = request(handler, "GET", "/api/v1/investors/2", "", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("get deleted = %d", w.Code)
	}
}

func TestBoundaryFailures(t *testing.T) {
	cases := []struct {
		name, method, path, body, contentType string
		status                                int
	}{
		{"malformed", "POST", "/api/v1/investors", `{`, "application/json", 400},
		{"empty body", "POST", "/api/v1/investors", ``, "application/json", 400},
		{"null", "POST", "/api/v1/investors", `null`, "application/json", 400},
		{"array", "POST", "/api/v1/investors", `[]`, "application/json", 400},
		{"unknown field", "POST", "/api/v1/investors", `{"id":1,"name":"Asha","country":"IN"}`, "application/json", 400},
		{"trailing value", "POST", "/api/v1/investors", `{"name":"Asha","country":"IN"} {}`, "application/json", 400},
		{"invalid domain", "POST", "/api/v1/investors", `{"name":" ","country":"IN"}`, "application/json", 400},
		{"oversized", "POST", "/api/v1/investors", `{"name":"` + strings.Repeat("x", 5000) + `","country":"IN"}`, "application/json", 413},
		{"trailing oversized whitespace", "POST", "/api/v1/investors", `{"name":"Asha","country":"IN"}` + strings.Repeat(" ", 5000), "application/json", 413},
		{"missing content type", "POST", "/api/v1/investors", `{}`, "", 415},
		{"wrong content type", "POST", "/api/v1/investors", `{}`, "text/plain", 415},
		{"invalid id", "GET", "/api/v1/investors/nope", "", "", 400},
		{"zero id", "GET", "/api/v1/investors/0", "", "", 400},
		{"missing item", "GET", "/api/v1/investors/99", "", "", 404},
		{"missing update", "PUT", "/api/v1/investors/99", `{"name":"Asha","country":"IN"}`, "application/json", 404},
		{"invalid update", "PUT", "/api/v1/investors/1", `{"name":"","country":"IN"}`, "application/json", 400},
		{"missing delete", "DELETE", "/api/v1/investors/99", "", "", 404},
		{"zero limit", "GET", "/api/v1/investors?limit=0", "", "", 400},
		{"large limit", "GET", "/api/v1/investors?limit=101", "", "", 400},
		{"invalid cursor", "GET", "/api/v1/investors?after=-1", "", "", 400},
		{"unknown query", "GET", "/api/v1/investors?page=1", "", "", 400},
		{"duplicate query", "GET", "/api/v1/investors?limit=1&limit=2", "", "", 400},
		{"malformed query", "GET", "/api/v1/investors?limit=%zz", "", "", 400},
		{"empty query value", "GET", "/api/v1/investors?after=", "", "", 400},
		{"wrong method", "PATCH", "/api/v1/investors/1", `{}`, "application/json", 405},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := newHandler(newStore())
			w := request(handler, tc.method, tc.path, tc.body, tc.contentType)
			if w.Code != tc.status {
				t.Fatalf("status = %d, want %d: %s", w.Code, tc.status, w.Body)
			}
			// Rejected writes must not create an
			w = request(handler, "GET", "/api/v1/investors", "", "")
			if strings.TrimSpace(w.Body.String()) != "[]" {
				t.Fatalf("invalid request mutated data: %s", w.Body)
			}
		})
	}
}

func TestHealthAndEmptyList(t *testing.T) {
	handler := newHandler(newStore())
	for _, path := range []string{"/health/live", "/api/v1/investors"} {
		w := request(handler, "GET", path, "", "")
		if w.Code != http.StatusOK || w.Header().Get("Content-Type") != "application/json" {
			t.Fatalf("%s = %d, %v", path, w.Code, w.Header())
		}
		body, err := io.ReadAll(w.Result().Body)
		if err != nil || !json.Valid(body) {
			t.Fatalf("invalid response body %q: %v", body, err)
		}
	}
}
