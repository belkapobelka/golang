package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusGetAllArticles(t *testing.T) {
	req, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllArticles)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("HandleFunc returned wrong ststus code: has %v expected %v", status, http.StatusOK)
	}

	expectedAnswer := `[{"Id":"1","Title":"First Title","Author":"First Author","Content":"First Content"},{"Id":"2","Title":"Second Title","Author":"Second Author","Content":"Second Content"}]`
	if recorder.Body.String() != expectedAnswer {
		t.Errorf("HandlerFunc returned answer: has %v expected %v", recorder.Body.String(), expectedAnswer)
	}
}

func TestGetArticleById(t *testing.T) {
	req, err := http.NewRequest("GET", "/article", nil)
	if err != nil {
		t.Fatal()
	}
	query := req.URL.Query()
	query.Add("id", "1")
	req.URL.RawQuery = query.Encode()

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(getArticleById)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("HandleFunc returned wrong ststus code: has %v expected %v", status, http.StatusOK)
	}

	expectedAnswer := `[{"Id":"1","Title":"First Title","Author":"First Author","Content":"First Content"}]`
	if recorder.Body.String() != expectedAnswer {
		t.Errorf("HandlerFunc returned answer: has %v expected %v", recorder.Body.String(), expectedAnswer)
	}
}
