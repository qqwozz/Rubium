// auth/provider.go
package auth

import "context"

// UserInfo — унифицированный профиль от любого провайдера
type UserInfo struct {
	ProviderID string // ID пользователя у провайдера
	FirstName  string
	LastName   string
	Email      string
	AvatarURL  string
}

// Provider — общий контракт для любого способа входа
type Provider interface {
	// Name возвращает идентификатор провайдера: "vk", "yandex"
	Name() string

	// AuthURL формирует ссылку для редиректа пользователя
	AuthURL(state, codeChallenge string) string

	// Exchange обменивает code на access_token
	Exchange(ctx context.Context, code, codeVerifier string) (accessToken string, err error)

	// UserInfo получает профиль пользователя по access_token
	UserInfo(ctx context.Context, accessToken string) (*UserInfo, error)
}