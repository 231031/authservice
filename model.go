package authservice

import "time"

type Auth struct {
	Id           string    `gorm:"id" json:"id"`
	Sub          string    `gorm:"sub" json:"sub"`
	RefreshToken string    `gorm:"refresh_token" json:"refresh_token"`
	Role         string    `gorm:"role" json:"role"`
	ExpiresIn    time.Time `gorm:"expires_in" json:"expires_in"`
}

type AuthPayload struct {
	Code        string `json:"code"`
	GrantType   string `json:"grant_type"`
	RedirectURI string `json:"redirect_uri"`
}

type AuthGoogleResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

type UserInfoResponse struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type TokenInfo struct {
	Id          string `json:"id"`
	AccessToken string `json:"access_token"`
	Role        string `json:"role"`
}

type RefreshTokenPayload struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}
