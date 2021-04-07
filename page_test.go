package spotify

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_NextPage(t *testing.T) {
	testTable := []struct {
		Name         string
		Input        *basePage
		ExpectedPath string
		Err          error
	}{
		{
			"success",
			&basePage{
				Next:  "/v1/albums/0sNOF9WDwhWunNAHPD3Baj/tracks",
				Total: 600,
			},
			"/v1/albums/0sNOF9WDwhWunNAHPD3Baj/tracks",
			nil,
		},
		{
			"no more pages",
			&basePage{
				Next: "",
			},
			"",
			ErrNoMorePages,
		},
		{
			"nil pointer error",
			nil,
			"",
			errors.New("spotify: p must be a non-nil pointer to a page"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			wasCalled := false
			client, server := testClientString(200, `{"total": 100}`, func(request *http.Request) {
				wasCalled = true
				assert.Equal(t, tt.ExpectedPath, request.URL.RequestURI())
			})
			if tt.Input != nil && tt.Input.Next != "" {
				tt.Input.Next = server.URL + tt.Input.Next // add fake server url so we intercept the message
			}

			err := client.NextPage(context.Background(), tt.Input)
			assert.Equal(t, tt.ExpectedPath != "", wasCalled)
			if tt.Err == nil {
				assert.NoError(t, err)
				assert.Equal(t, 100, tt.Input.Total) // value should be from original 600
			} else {
				assert.EqualError(t, err, tt.Err.Error())
			}
		})
	}
}

func TestClient_PreviousPage(t *testing.T) {
	testTable := []struct {
		Name         string
		Input        *basePage
		ExpectedPath string
		Err          error
	}{
		{
			"success",
			&basePage{
				Previous: "/v1/albums/0sNOF9WDwhWunNAHPD3Baj/tracks",
				Total:    600,
			},
			"/v1/albums/0sNOF9WDwhWunNAHPD3Baj/tracks",
			nil,
		},
		{
			"no more pages",
			&basePage{
				Previous: "",
			},
			"",
			ErrNoMorePages,
		},
		{
			"nil pointer error",
			nil,
			"",
			errors.New("spotify: p must be a non-nil pointer to a page"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			wasCalled := false
			client, server := testClientString(200, `{"total": 100}`, func(request *http.Request) {
				wasCalled = true
				assert.Equal(t, tt.ExpectedPath, request.URL.RequestURI())
			})
			if tt.Input != nil && tt.Input.Previous != "" {
				tt.Input.Previous = server.URL + tt.Input.Previous // add fake server url so we intercept the message
			}

			err := client.PreviousPage(context.Background(), tt.Input)
			assert.Equal(t, tt.ExpectedPath != "", wasCalled)
			if tt.Err == nil {
				assert.NoError(t, err)
				assert.Equal(t, 100, tt.Input.Total) // value should be from original 600
			} else {
				assert.EqualError(t, err, tt.Err.Error())
			}
		})
	}
}
