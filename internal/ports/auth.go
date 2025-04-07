package ports

type AuthDto struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresIn  int    `json:"access_token_expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	TokenType             string `json:"token_type"`
}

type Authenticator interface {
	Authenticate(username, password string) (AuthDto, error)
	RefreshToken(token string) (AuthDto, error)
}
