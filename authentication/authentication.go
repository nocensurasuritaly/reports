package authentication

import (
	"context"
	"github.com/EgidioCaprino/reddit-oauth2/token"
	"log"
	"net/http"
)

const ContextKeyWebToken = "ContextKeyWebToken"

var decryptWebToken = token.DecryptWebToken

func WithAuthentication(decryptionKey string, next http.HandlerFunc) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		authorization := request.Header.Get("Authorization")
		if len(authorization) == 0 {
			log.Println("request is missing authorization header")
			responseWriter.WriteHeader(http.StatusForbidden)
			return
		}
		webToken, err := decryptWebToken(authorization, decryptionKey)
		if err != nil {
			log.Println("unable to decrypt web token:", err)
			responseWriter.WriteHeader(http.StatusForbidden)
			return
		}
		newContext := context.WithValue(request.Context(), ContextKeyWebToken, webToken)
		newRequest := request.WithContext(newContext)
		next(responseWriter, newRequest)
	}
}
