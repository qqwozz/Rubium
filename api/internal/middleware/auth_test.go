package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/supabase"

	"github.com/gin-gonic/gin"
)

func mockSupabaseAuthServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/v1/user":
			token := r.Header.Get("Authorization")
			if token == "Bearer valid-token" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id":"user-123"}`))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"invalid token"}`))
			}
		case "/rest/v1/rubium_users":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id":"rubium-123"}]`))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		}
	}))
}

func TestRequireAuth_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	srv := mockSupabaseAuthServer()
	defer srv.Close()

	client := supabase.NewClient(srv.URL, "anon", "service")
	r := gin.New()
	r.Use(RequireAuth(client))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"user_id": c.GetString("user_id")})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRequireAuth_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	srv := mockSupabaseAuthServer()
	defer srv.Close()

	client := supabase.NewClient(srv.URL, "anon", "service")
	r := gin.New()
	r.Use(RequireAuth(client))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestOptionalAuth_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	srv := mockSupabaseAuthServer()
	defer srv.Close()

	client := supabase.NewClient(srv.URL, "anon", "service")
	r := gin.New()
	r.Use(OptionalAuth(client))
	r.GET("/public", func(c *gin.Context) {
		_, exists := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{"authed": exists})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/public", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		header   string
		expected string
	}{
		{"Bearer token123", "token123"},
		{"Basic token123", ""},
		{"", ""},
		{"Bearer ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest("GET", "/", nil)
			c.Request.Header.Set("Authorization", tt.header)
			got := extractToken(c)
			if got != tt.expected {
				t.Errorf("extractToken(%q) = %q, want %q", tt.header, got, tt.expected)
			}
		})
	}
}