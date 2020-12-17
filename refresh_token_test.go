package moneybutton

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockHTTPRefreshAccessToken for mocking requests
type mockHTTPRefreshAccessToken struct{}

// Do is a mock http request
func (m *mockHTTPRefreshAccessToken) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	if req.URL.String() == endpointToken {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"access_token":"eyJ0zXAiOiJKV1QzLCJhbGciOzJIUzI1NiJ9.eyJzdWIiOiIyNjQ1IiwiYXVkIjoiMjhzNzUyZDg0NmRmMjU0YzM3ZWQ1MjBlZTlkYjJjMzMiLCJleHAiOjE2MDzyMzc5MjIsInNjb3BlIjoidXNlcnMucHJvZmlsZXM6cmVhZCBhdXRoLnVzZXJfaWRlbzRpdHk6cmVhZCJ9.bTfOzUvLE5zd2IvRQMXVCJX2kzOB9-44EGLn92UzAMI","token_type":"` + authHeaderBearer + `","expires_in":3600,"scope":"` + PermissionsIdentity + " " + PermissionsProfile + `","refresh_token":"379fc72za81ae1ze2f958399bz2c990350f46034z840584cc5dec63879b8c876"}`)))
	}

	// Default is valid
	return resp, nil
}

func TestClient_RefreshAccessToken(t *testing.T) {
	t.Parallel()

	t.Run("missing all parameters", func(t *testing.T) {
		client := newTestClient(&mockHTTPRefreshAccessToken{})
		assert.NotNil(t, client)
		tokenResponse, err := client.RefreshAccessToken(
			context.Background(),
			"",
			"",
		)
		assert.Error(t, err)
		assert.Nil(t, tokenResponse)
	})

	t.Run("missing client id", func(t *testing.T) {
		client := newTestClient(&mockHTTPRefreshAccessToken{})
		assert.NotNil(t, client)
		tokenResponse, err := client.RefreshAccessToken(
			context.Background(),
			"",
			"1234567",
		)
		assert.Error(t, err)
		assert.Nil(t, tokenResponse)
	})

	t.Run("missing access token", func(t *testing.T) {
		client := newTestClient(&mockHTTPRefreshAccessToken{})
		assert.NotNil(t, client)
		tokenResponse, err := client.RefreshAccessToken(
			context.Background(),
			"1234567",
			"",
		)
		assert.Error(t, err)
		assert.Nil(t, tokenResponse)
	})

	t.Run("api error response", func(t *testing.T) {
		client := newTestClient(&mockHTTPAPIError{})
		assert.NotNil(t, client)
		tokenResponse, err := client.RefreshAccessToken(
			context.Background(),
			"1234567",
			"1234567",
		)
		assert.Error(t, err)
		assert.Nil(t, tokenResponse)
	})

	t.Run("http error", func(t *testing.T) {
		client := newTestClient(&mockHTTPError{})
		assert.NotNil(t, client)
		tokenResponse, err := client.RefreshAccessToken(
			context.Background(),
			"1234567",
			"1234567",
		)
		assert.Error(t, err)
		assert.Nil(t, tokenResponse)
	})

	t.Run("valid response", func(t *testing.T) {
		client := newTestClient(&mockHTTPRefreshAccessToken{})
		assert.NotNil(t, client)
		tokenResponse, err := client.RefreshAccessToken(
			context.Background(),
			"1234567",
			"1234567",
		)
		assert.NoError(t, err)
		assert.NotNil(t, tokenResponse)
		assert.Equal(t, "eyJ0zXAiOiJKV1QzLCJhbGciOzJIUzI1NiJ9.eyJzdWIiOiIyNjQ1IiwiYXVkIjoiMjhzNzUyZDg0NmRmMjU0YzM3ZWQ1MjBlZTlkYjJjMzMiLCJleHAiOjE2MDzyMzc5MjIsInNjb3BlIjoidXNlcnMucHJvZmlsZXM6cmVhZCBhdXRoLnVzZXJfaWRlbzRpdHk6cmVhZCJ9.bTfOzUvLE5zd2IvRQMXVCJX2kzOB9-44EGLn92UzAMI", tokenResponse.AccessToken)
		assert.Equal(t, uint32(3600), tokenResponse.ExpiresIn)
		assert.Equal(t, "379fc72za81ae1ze2f958399bz2c990350f46034z840584cc5dec63879b8c876", tokenResponse.RefreshToken)
		assert.Equal(t, PermissionsIdentity+" "+PermissionsProfile, tokenResponse.Scope)
		assert.Equal(t, authHeaderBearer, tokenResponse.TokenType)
	})
}
