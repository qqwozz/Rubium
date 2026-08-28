package supabase

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockPostgrestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			if r.URL.Path == "/rest/v1/tasks" {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]map[string]interface{}{
					{"id": "task-1", "title": "Test"},
				})
				return
			}
		case "POST":
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{"id": "new-id"})
			return
		case "PATCH":
			w.WriteHeader(http.StatusNoContent)
			return
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
}

func TestClient_Query(t *testing.T) {
	srv := mockPostgrestServer()
	defer srv.Close()

	client := NewClient(srv.URL, "anon", "service")
	var tasks []map[string]interface{}
	err := client.Query(context.Background(), "tasks", false, &tasks)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
}

func TestClient_Post(t *testing.T) {
	srv := mockPostgrestServer()
	defer srv.Close()

	client := NewClient(srv.URL, "anon", "service")
	var result map[string]interface{}
	err := client.Post(context.Background(), "tasks", true, map[string]interface{}{"title": "New"}, &result)
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	if result["id"] != "new-id" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestClient_Patch(t *testing.T) {
	srv := mockPostgrestServer()
	defer srv.Close()

	client := NewClient(srv.URL, "anon", "service")
	err := client.Patch(context.Background(), "tasks?id=eq.1", true, map[string]interface{}{"title": "Updated"})
	if err != nil {
		t.Fatalf("Patch failed: %v", err)
	}
}

func TestClient_Delete(t *testing.T) {
	srv := mockPostgrestServer()
	defer srv.Close()

	client := NewClient(srv.URL, "anon", "service")
	err := client.Delete(context.Background(), "tasks?id=eq.1", true)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestClient_RPC(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rest/v1/rpc/my_func" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"result": 42}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "anon", "service")
	var result map[string]interface{}
	err := client.RPC(context.Background(), "my_func", true, map[string]interface{}{}, &result)
	if err != nil {
		t.Fatalf("RPC failed: %v", err)
	}
}

func TestClient_AuthUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/v1/user" {
			token := r.Header.Get("Authorization")
			if token == "Bearer good-token" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id":"user-1"}`))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "anon", "service")
	uid, err := client.AuthUser(context.Background(), "good-token")
	if err != nil {
		t.Fatalf("AuthUser failed: %v", err)
	}
	if uid != "user-1" {
		t.Errorf("expected user-1, got %s", uid)
	}

	_, err = client.AuthUser(context.Background(), "bad-token")
	if err == nil {
		t.Error("expected error for bad token")
	}
}

func TestClient_RawQuery_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"db error"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "anon", "service")
	_, err := client.RawQuery(context.Background(), "tasks", false)
	if err == nil {
		t.Error("expected error for 500 response")
	}
}