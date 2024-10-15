package token

import (
	"time"

	"gihub.com/saiddis/quizgo/internal/install/database"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserClaims struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	jwt.StandardClaims
}

func NewUserClaims(u database.User) *UserClaims {

	return &UserClaims{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
}
