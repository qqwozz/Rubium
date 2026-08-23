package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/supabase"

	"github.com/gin-gonic/gin"
)

func mockNotebooksServer(t *testing.T) (*httptest.Server, *supabase.Client) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/auth/v1/user":
			token := r.Header.Get("Authorization")
			if token == "Bearer valid-token" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id":"auth-123"}`))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		case r.URL.Path == "/rest/v1/rubium_users":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id":"rubium-123"}]`))
		case r.URL.Path == "/rest/v1/notebooks" && r.Method == "GET":
			if r.URL.RawQuery == "select=*&id=eq.valid-uuid&limit=1" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[{"id":"valid-uuid","user_id":"rubium-123","title":"Test","is_public":true,"average_rating":0,"ratings_count":0,"views_count":0,"copies_count":0,"created_at":"2024-01-01","updated_at":"2024-01-01","content":null}]`))
			} else if r.URL.RawQuery == "select=user_id&id=eq.valid-uuid&limit=1" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[{"user_id":"rubium-123"}]`))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			}
		case r.URL.Path == "/rest/v1/notebooks" && r.Method == "POST":
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`[{"id":"new-uuid","user_id":"rubium-123","title":"New Notebook","is_public":false,"average_rating":0,"ratings_count":0,"views_count":0,"copies_count":0,"created_at":"2024-01-01","updated_at":"2024-01-01","content":null}]`))
		case r.URL.Path == "/rest/v1/notebooks" && r.Method == "PATCH":
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/rest/v1/notebooks" && r.Method == "DELETE":
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/rest/v1/rpc/increment_notebook_views":
			w.WriteHeader(http.StatusOK)
		case r.URL.Path == "/rest/v1/rpc/increment_notebook_copies":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		}
	}))
	return srv, supabase.NewClient(srv.URL, "anon", "service")
}

func setupNotebooksRouter(t *testing.T) (*gin.Engine, *httptest.Server) {
	gin.SetMode(gin.TestMode)
	srv, client := mockNotebooksServer(t)
	h := NewNotebooksHandler(client)

	r := gin.New()
	r.GET("/api/v1/notebooks/community", h.GetCommunityNotebooks)
	r.GET("/api/v1/notebooks/:id", h.GetNotebookByID)
	r.GET("/api/v1/notebooks/:id/rating", h.GetRating)

	auth := r.Group("/api/v1/notebooks")
	auth.Use(func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "Bearer valid-token" {
			c.Set("user_id", "auth-123")
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
		}
	})
	{
		auth.GET("", h.GetNotebooks)
		auth.POST("", h.CreateNotebook)
		auth.PUT("/:id", h.UpdateNotebook)
		auth.DELETE("/:id", h.DeleteNotebook)
		auth.POST("/:id/copy", h.CopyNotebook)
		auth.POST("/:id/rate", h.RateNotebook)
		auth.POST("/:id/view", h.IncrementViews)
	}

	return r, srv
}

func TestGetNotebookByID_Public(t *testing.T) {
	r, srv := setupNotebooksRouter(t)
	defer srv.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notebooks/valid-uuid", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetNotebookByID_InvalidUUID(t *testing.T) {
	r, srv := setupNotebooksRouter(t)
	defer srv.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notebooks/bad-uuid", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateNotebook_Unauthorized(t *testing.T) {
	r, srv := setupNotebooksRouter(t)
	defer srv.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notebooks", bytes.NewBuffer([]byte(`{"title":"Test"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestCreateNotebook_Success(t *testing.T) {
	r, srv := setupNotebooksRouter(t)
	defer srv.Close()

	body, _ := json.Marshal(map[string]string{"title": "Test Notebook"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notebooks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestIncrementViews_InvalidUUID(t *testing.T) {
	r, srv := setupNotebooksRouter(t)
	defer srv.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notebooks/bad-uuid/view", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetCommunityNotebooks(t *testing.T) {
	r, srv := setupNotebooksRouter(t)
	defer srv.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notebooks/community", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if _, ok := resp["notebooks"]; !ok {
		t.Error("expected notebooks in response")
	}
}