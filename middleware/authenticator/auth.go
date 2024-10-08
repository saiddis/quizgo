package authenticator

import (
	"context"
	"fmt"
	"log"
	"time"

	"gihub.com/saiddis/quizgo/internal/install/database"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	db       *database.Queries
}

func New(db *database.Queries) (*Authenticator, error) {
	env, err := godotenv.Read("../../.env")
	if err != nil {
		log.Fatalf("Could't read .env file: %v", err)
	}
	secret := env["AUTH0_CLIENT_SECRET"]
	domain := env["AUTH0_DOMAIN"]
	clientID := env["AUTH0_CLIENT_ID"]
	callbackURL := env["AUTH0_CALLBACK_URL"]

	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+domain+"/",
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		RedirectURL:  callbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "user_id"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		db:       db,
	}, nil
}

func (a *Authenticator) VerifyIDToken(c *gin.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)

	if !ok {
		return nil, fmt.Errorf("no token_id field in oath2 token")
	}

	userInfo, err := a.Provider.UserInfo(c, oauth2.StaticTokenSource(token))
	if err != nil {
		return nil, fmt.Errorf("couln't get user info: %v", err)
	}

	user, err := a.db.GetUserByEmail(c, userInfo.Email)
	if err != nil {
		user, err = a.db.CreateUser(c, database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			Email:     userInfo.Email,
		})
		err = userInfo.Claims(user.ID)
		if err != nil {
			log.Println(err.Error())
		}
	}

	oidcConfig := &oidc.Config{
		ClientID: a.Config.ClientID,
	}

	return a.Provider.Verifier(oidcConfig).Verify(c, rawIDToken)

}
