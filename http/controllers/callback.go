package controllers

import (
	"log"
	"net/http"

	"gihub.com/saiddis/quizgo/middleware/authenticator"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CallbackHandler(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if c.Query("state") != session.Get("state") {
			c.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := auth.Config.Exchange(c.Request.Context(), c.Query("code"))
		if err != nil {
			c.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
			return
		}

		idToken, err := auth.VerifyIDToken(c, token)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			log.Printf("failed to verify id token: %v", token)
			log.Printf("session: %v", session)
			log.Printf("idToken: %v", idToken)
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		if err := session.Save(); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Redirect to logged in page.
		c.Redirect(http.StatusTemporaryRedirect, "/user")
	}
}
