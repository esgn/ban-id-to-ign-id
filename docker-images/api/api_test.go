package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestEmptyPosition(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BanToIgn)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Statut HTTP invalide: attendu %v obtenu %v",
			status, http.StatusNotFound)
	}
	expected := `{"error":{"message":"Veuillez indiquer au moins une cle_interop"}}`
	if rr.Body.String() != expected {
		t.Errorf("Réponse invalide : attendu %v obtenu %v",
			rr.Body.String(), expected)
	}
}

func TestInvalidCleInterop(t *testing.T) {
	req, err := http.NewRequest("GET", "xxx", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BanToIgn)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Statut HTTP invalide: attendu %v obtenu %v",
			status, http.StatusNotFound)
	}
	expected := `{"error":{"message":"xxx est une cle_interop invalide"}}`
	if rr.Body.String() != expected {
		t.Errorf("Réponse invalide : attendu %v obtenu %v",
			rr.Body.String(), expected)
	}
}

func TestCommas(t *testing.T) {
	req, err := http.NewRequest("GET", ",,,,,", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BanToIgn)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Statut HTTP invalide: attendu %v obtenu %v",
			status, http.StatusNotFound)
	}
	expected := `{"error":{"message":"Veuillez indiquer au moins une cle_interop"}}`
	if rr.Body.String() != expected {
		t.Errorf("Réponse invalide : attendu %v obtenu %v",
			rr.Body.String(), expected)
	}
}

func TestTooManyCleInterop(t *testing.T) {

	path := ""
	for i := 0; i < (maxIds + 1); i++ {
		path += strconv.Itoa(i) + ","
	}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BanToIgn)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Statut HTTP invalide: attendu %v obtenu %v",
			status, http.StatusNotFound)
	}
	expected := `{"error":{"message":"Liste de cle_interop dépassant la limite de ` + strconv.Itoa(maxIds) + `"}}`
	if rr.Body.String() != expected {
		t.Errorf("Réponse invalide : attendu %v obtenu %v",
			rr.Body.String(), expected)
	}
}
