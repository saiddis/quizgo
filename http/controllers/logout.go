package controllers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Handler for our logout.
func LogoutHandler(c *gin.Context) {

	env, err := godotenv.Read("../../.env")
	if err != nil {
		log.Fatalf("Could't read .env file: %v", err)
	}
	domain := env["AUTH0_DOMAIN"]
	clientID := env["AUTH0_CLIENT_ID"]
	logoutUrl, err := url.Parse("https://" + domain + "/v2/logout")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + c.Request.Host)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", clientID)
	logoutUrl.RawQuery = parameters.Encode()

	c.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}
