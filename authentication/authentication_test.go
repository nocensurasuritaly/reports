package authentication

import (
	"errors"
	"github.com/EgidioCaprino/reddit-oauth2/token"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithAuthenticationFailsIfAuthorizationHeaderIsMissing(t *testing.T) {
	key := "key"
	next := func(responseWriter http.ResponseWriter, request *http.Request) {}
	handler := WithAuthentication(key, next)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	handler.ServeHTTP(responseRecorder, request)
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}

func TestWithAuthenticationDecryptsTheWebToken(t *testing.T) {
	var actualEncryptedWebToken string
	var actualKey string
	key := "key"
	next := func(responseWriter http.ResponseWriter, request *http.Request) {}
	decryptWebToken = func(encryptedWebToken string, key string) (*token.WebToken, error) {
		actualEncryptedWebToken = encryptedWebToken
		actualKey = key
		return &token.WebToken{}, nil
	}
	handler := WithAuthentication(key, next)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.Header.Set("Authorization", "encrypted web token")
	handler.ServeHTTP(responseRecorder, request)
	assert.Equal(t, "encrypted web token", actualEncryptedWebToken)
	assert.Equal(t, "key", actualKey)
}

func TestWithAuthenticationFailsWhenCannotDecryptTheWebToken(t *testing.T) {
	key := "key"
	next := func(responseWriter http.ResponseWriter, request *http.Request) {}
	decryptWebToken = func(encryptedWebToken string, key string) (*token.WebToken, error) {
		return nil, errors.New("")
	}
	handler := WithAuthentication(key, next)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.Header.Set("Authorization", "encrypted web token")
	handler.ServeHTTP(responseRecorder, request)
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}

func TestWithAuthenticationInvokeNextHandler(t *testing.T) {
	invoked := false
	key := "key"
	next := func(responseWriter http.ResponseWriter, request *http.Request) {
		invoked = true
	}
	decryptWebToken = func(encryptedWebToken string, key string) (*token.WebToken, error) {
		return &token.WebToken{}, nil
	}
	handler := WithAuthentication(key, next)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.Header.Set("Authorization", "encrypted web token")
	handler.ServeHTTP(responseRecorder, request)
	assert.True(t, invoked)
}

func TestWithAuthenticationAddsWebTokenToTheContext(t *testing.T) {
	expected := &token.WebToken{}
	var actual *token.WebToken
	key := "key"
	next := func(responseWriter http.ResponseWriter, request *http.Request) {
		actual = request.Context().Value(ContextKeyWebToken).(*token.WebToken)
	}
	decryptWebToken = func(encryptedWebToken string, key string) (*token.WebToken, error) {
		return expected, nil
	}
	handler := WithAuthentication(key, next)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.Header.Set("Authorization", "encrypted web token")
	handler.ServeHTTP(responseRecorder, request)
	assert.Same(t, expected, actual)
}
