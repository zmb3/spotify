package spotify

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func testClient(code int, body io.Reader, validators ...func(*http.Request)) (*Client, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range validators {
			v(r)
		}
		w.WriteHeader(code)
		_, _ = io.Copy(w, body)
		r.Body.Close()
		if closer, ok := body.(io.Closer); ok {
			closer.Close()
		}
	}))
	client := &Client{
		http:    http.DefaultClient,
		baseURL: server.URL + "/",
	}
	return client, server
}

// Returns a client whose requests will always return
// the specified status code and body.
func testClientString(code int, body string, validators ...func(*http.Request)) (*Client, *httptest.Server) {
	return testClient(code, strings.NewReader(body), validators...)
}

// Returns a client whose requests will always return
// a response with the specified status code and a body
// that is read from the specified file.
func testClientFile(code int, filename string, validators ...func(*http.Request)) (*Client, *httptest.Server) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return testClient(code, f, validators...)
}

func TestClient_Token(t *testing.T) {
	// oauth setup for valid test token
	config := oauth2.Config{
		ClientID:     "test_client",
		ClientSecret: "test_client_secret",
		Endpoint:     oauth2.Endpoint{},
		RedirectURL:  "http://redirect.redirect",
		Scopes:       nil,
	}
	token := &oauth2.Token{
		AccessToken:  "access_token",
		TokenType:    "test_type",
		RefreshToken: "refresh_token",
		Expiry:       time.Now().Add(time.Hour),
	}

	t.Run("success", func(t *testing.T) {
		httpClient := config.Client(context.Background(), token)
		client := New(httpClient)
		token, err := client.Token()
		if err != nil {
			t.Error(err)
		}

		if !token.Valid() {
			t.Errorf("Token should be valid: %v", token)
		}
		if token.AccessToken != "access_token" {
			t.Errorf("Invalid access token data: %s", token.AccessToken)
		}
		if token.RefreshToken != "refresh_token" {
			t.Errorf("Invalid refresh token data: %s", token.RefreshToken)
		}
		if token.TokenType != "test_type" {
			t.Errorf("Invalid token type: %s", token.TokenType)
		}
	})

	t.Run("non oauth2 transport", func(t *testing.T) {
		client := &Client{
			http: http.DefaultClient,
		}
		_, err := client.Token()
		if err == nil || err.Error() != "spotify: client not backed by oauth2 transport" {
			t.Errorf("Should throw error: %s", "spotify: client not backed by oauth2 transport")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		httpClient := config.Client(context.Background(), nil)
		client := New(httpClient)
		_, err := client.Token()
		if err == nil || err.Error() != "oauth2: token expired and refresh token is not set" {
			t.Errorf("Should throw error: %s", "oauth2: token expired and refresh token is not set")
		}
	})
}

func TestDecode429Error(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Retry-After": []string{"2"}},
		Body:       io.NopCloser(strings.NewReader(`Too many requests`)),
	}

	err := decodeError(resp)
	if err == nil {
		t.Fatal("Expected error")
	}
	if err.Error() != "spotify: Too many requests [429]" {
		t.Error("Unexpected error message:", err.Error())
	}
	const wantSTatus = http.StatusTooManyRequests
	var gotStatus int
	var statusErr interface{ HTTPStatus() int }
	if errors.As(err, &statusErr) {
		gotStatus = statusErr.HTTPStatus()
	}
	if gotStatus != wantSTatus {
		t.Errorf("Expected status %d, got %d", wantSTatus, gotStatus)
	}
}
