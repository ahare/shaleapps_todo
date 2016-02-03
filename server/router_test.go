package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	r := &router{}

	r.HandleFunc("/", "GET", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "index GET")
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Body.String() != "index GET" {
		t.Fatal("Index call failed")
	}

	r.HandleFunc("/", "POST", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "index POST")
	})

	req, _ = http.NewRequest("POST", "/", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Body.String() != "index POST" {
		t.Fatal("Index call failed")
	}
}
