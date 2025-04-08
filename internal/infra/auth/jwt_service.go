package auth

import (
	"time"	

	"github.com/golang-jwt/jwt/v4"
	"errors"
	"github.com/PedroNetto404/marmitech-backend/internal/domain/ports"
)

type jwtService struct {
	secretKey string
	issuer string
	audience string
	expirationMinutes int
}

func NewJwtService(secretKey, issuer, audience string, expirationMinutes int) ports.IAuthService {
	return &jwtService{
		secretKey: secretKey,
		issuer: issuer,
		audience: audience,
		expirationMinutes: expirationMinutes,
	}
}

func (s *jwtService) Generate(payload ports.AuthPayload) (ports.AuthTokens, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": payload.Email,
		"pwd": payload.Pwd,
		"exp": time.Now().Add(time.Duration(s.expirationMinutes) * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return ports.AuthTokens{}, err
	}

	return ports.AuthTokens{
		AccessToken: tokenString,
		AccessTokenExpiresAt: time.Now().Add(time.Duration(s.expirationMinutes) * time.Minute),
	}, nil
}

func (s *jwtService) Validate(token string) (ports.AuthPayload, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return ports.AuthPayload{}, err
	}	

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return ports.AuthPayload{}, errors.New("invalid token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return ports.AuthPayload{}, errors.New("email not found in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return ports.AuthPayload{}, errors.New("role not found in token")
	}
		
	return ports.AuthPayload{Email: email, Role: role}, nil
}


