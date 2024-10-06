package token

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type Service struct {
	secret string
}

func NewService(secret string) *Service {
	return &Service{
		secret: secret,
	}
}

func (f *Service) NewAccess(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(f.secret))
}

func (f *Service) NewRefresh(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(f.secret))
}

func (f *Service) ParseAccess(accessToken string) (*UserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(f.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error parsing access token: %v", err)
	}

	return parsedAccessToken.Claims.(*UserClaims), nil
}

func (f *Service) ParseRefresh(refreshToken string, secret string) (*jwt.StandardClaims, error) {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(f.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error parsing refresh token: %v", err)
	}
	return parsedRefreshToken.Claims.(*jwt.StandardClaims), nil
}
