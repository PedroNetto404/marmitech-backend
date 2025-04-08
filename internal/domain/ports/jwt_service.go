package ports

import "time"

type (
	AuthTokens struct {
		AccessToken string `json:"access_token"`
		AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	}

	AuthPayload struct {
		Email string `json:"email"`
		Pwd string `json:"pwd"`
		Role string `json:"role"`
	}

	IAuthService interface {
		Generate(payload AuthPayload) (AuthTokens, error)
		Validate(token string) (AuthPayload, error)
	}
)
