package authservice

type User struct {
	id            string
	sub           string
	refresh_token string
	role          string
}

type GoogleInfo struct {
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
