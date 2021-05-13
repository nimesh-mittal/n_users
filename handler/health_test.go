package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHome(t *testing.T) {

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8085/", nil)

	NewHealthHandler().NewHealthRouter().ServeHTTP(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("home didn’t respond 200 OK: %s", resp.Status)
	}

	type home struct {
		Greet string
	}
	var sr home
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		t.Errorf("home response parsing error %s", err)
	}

	if sr.Greet != "hello" {
		t.Errorf("home response status is %s but expected true", sr.Greet)
	}
}

func TestHealth(t *testing.T) {

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8085/_health", nil)

	NewHealthHandler().NewHealthRouter().ServeHTTP(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("health didn’t respond 200 OK: %s", resp.Status)
	}

	type health struct {
		Status string
	}
	var sr health
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		t.Errorf("health response parsing error %s", err)
	}

	if sr.Status != "green" {
		t.Errorf("health response status is %s but expected true", sr.Status)
	}
}
