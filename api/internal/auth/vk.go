// auth/vk.go
package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type VKProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scope        string
}

func (p *VKProvider) Name() string { return "vk" }

func (p *VKProvider) AuthURL(state, codeChallenge string) string {
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", p.ClientID)
	params.Set("redirect_uri", p.RedirectURI)
	params.Set("scope", p.Scope)
	params.Set("state", state)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")
	return "https://id.vk.ru/authorize?" + params.Encode()
}

func (p *VKProvider) Exchange(ctx context.Context, code, codeVerifier string) (string, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("code_verifier", codeVerifier)
	form.Set("client_id", p.ClientID)
	form.Set("redirect_uri", p.RedirectURI)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://id.vk.ru/oauth2/auth",
		bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vk exchange failed: %s", resp.Status)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

func (p *VKProvider) UserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	form := url.Values{}
	form.Set("access_token", accessToken)
	form.Set("client_id", p.ClientID)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://id.vk.ru/oauth2/user_info",
		bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw struct {
		User struct {
			UserID    string `json:"user_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
			Avatar    string `json:"avatar"`
		} `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return &UserInfo{
		ProviderID: raw.User.UserID,
		FirstName:  raw.User.FirstName,
		LastName:   raw.User.LastName,
		Email:      raw.User.Email,
		AvatarURL:  raw.User.Avatar,
	}, nil
}