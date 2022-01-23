package webhook

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewDeliverer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
		},
	))
	defer server.Close()
	conf := Config{
		Target: server.URL,
	}

	_, err := New(&conf)
	if err == nil {
		t.Fatalf("Expected error but error is nil")
	}
	conf.Client = server.Client()
	_, err = New(&conf)
	if err != nil {
		t.Fatalf("New failed:: %v", err)
	}
}

func TestDeliver(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
		},
	))
	defer server.Close()
	conf := Config{
		Client: server.Client(),
		Target: server.URL,
	}

	d, err := New(&conf)
	if err != nil {
		t.Fatalf("New failed:: %v", err)
	}
	err = d.Deliver(context.Background(), "")
	if err != nil {
		t.Fatalf("New failed:: %v", err)
	}

}
