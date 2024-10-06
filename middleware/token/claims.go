package token

import (
	"time"

	"gihub.com/saiddis/quizgo/internal/install/database"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type UserClaims struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	jwt.StandardClaims
}

func NewUserClaims(u database.User) *UserClaims {

	return &UserClaims{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
}
