package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/oauth2/v2"
)

var httpClient = &http.Client{}

// VerifyAPIRequest takes a request and validates the token
func VerifyAPIRequest(c *gin.Context, tokens []string) bool {
	if tokens == nil || len(tokens) == 0 {
		c.JSON(400, "Missing ID token")
		return false
	}
	idToken := tokens[0]
	_, err := VerifyIDToken(idToken)
	if err != nil {
		c.JSON(400, "Invalid ID token")
		return false
	}

	return true
}

// VerifyIDToken verifies the session ID token to approve API requests
func VerifyIDToken(idToken string) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.New(httpClient)
	if err != nil {
		return nil, err
	}
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
}
