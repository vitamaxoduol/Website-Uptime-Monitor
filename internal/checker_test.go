package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckWebsiteOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
		},
	))
	defer server.Close()

	result, err := CheckWebsite(server.Client(), server.URL)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", result.StatusCode)
	}

	if result.URL != server.URL {
		t.Errorf("expected status %q, got %q", server.URL, result.URL)
	}

}

func TestCheckWebsiteServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusInternalServerError)
		},
	))
	defer server.Close()

	result, err := CheckWebsite(server.Client(), server.URL)

	if err == nil {
		t.Fatalf("expected an error for status 500, got nil")
	}

	if result.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected an error status 500, got %d", result.StatusCode)
	}
}

func TestCheckWebsiteInvalidURL(t *testing.T) {
	client := &http.Client{}
	url := "://invalid"

	result, err := CheckWebsite(client, url)

	if err == nil {
		t.Fatalf("expected an error for an invalid URL, got nil")
	}

	if result.URL != url {
		t.Errorf("expected URL %q, got %q", url, result.URL)
	}

	if result.StatusCode != 0 {
		t.Errorf("expected status 0 without a reponse, got %d", result.StatusCode)
	}
}
