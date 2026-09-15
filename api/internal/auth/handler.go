// auth/handler.go
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"context"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	registry *Registry
	// storage — интерфейс к вашему хранилищу (БД, Redis)
	storage Storage
}

// Storage — интерфейс для сохранения PKCE/state и токенов.
// Реализуйте под свою БД.
type Storage interface {
	SaveAuthSession(ctx context.Context, state, provider, codeVerifier string) error
	GetAuthSession(ctx context.Context, state string) (provider, codeVerifier string, err error)
	DeleteAuthSession(ctx context.Context, state string) error

	// FindOrCreateUser — находит или создаёт пользователя по данным провайдера
	FindOrCreateUser(ctx context.Context, provider string, info *UserInfo) (userID string, err error)
}

// GET /auth/:provider/url
func (h *Handler) GetAuthURL(c *gin.Context) {
	providerName := c.Param("provider")

	provider, err := h.registry.Get(providerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// PKCE
	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)
	state := generateState()

	if err := h.storage.SaveAuthSession(c.Request.Context(), state, providerName, codeVerifier); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить сессию"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": provider.AuthURL(state, codeChallenge),
	})
}

// POST /auth/:provider/callback
type CallbackRequest struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}

func (h *Handler) Callback(c *gin.Context) {
	providerName := c.Param("provider")

	provider, err := h.registry.Get(providerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req CallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "нужны code и state"})
		return
	}

	// Проверяем state
	storedProvider, codeVerifier, err := h.storage.GetAuthSession(c.Request.Context(), req.State)
	if err != nil || storedProvider != providerName {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный state"})
		return
	}

	// Обмен кода на токен
	accessToken, err := provider.Exchange(c.Request.Context(), req.Code, codeVerifier)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "не удалось получить токен"})
		return
	}

	// Профиль
	info, err := provider.UserInfo(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "не удалось получить профиль"})
		return
	}

	// Find-or-create локального пользователя
	userID, err := h.storage.FindOrCreateUser(c.Request.Context(), providerName, info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания пользователя"})
		return
	}

	_ = h.storage.DeleteAuthSession(c.Request.Context(), req.State)

	// Выдайте свою сессию/JWT
	c.JSON(http.StatusOK, gin.H{
		"user_id":    userID,
		"first_name": info.FirstName,
		"last_name":  info.LastName,
		"email":      info.Email,
		"provider":   providerName,
	})
}

// GET /auth/providers — список доступных провайдеров
func (h *Handler) ListProviders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"providers": h.registry.List()})
}

// --- утилиты ---

func generateCodeVerifier() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func generateState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}