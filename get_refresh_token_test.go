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

// mockHTTPGetRefreshToken for mocking requests
type mockHTTPGetRefreshToken struct{}

// Do is a mock http request
func (m *mockHTTPGetRefreshToken) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	if req.URL.String() == endpointToken {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"access_token":"zzz0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIzOiIyNjQ1IiwiYXVkIjoiMjhmNzUyZDg0NmRmzjU0YzM3ZWQ1MjBlZTlkYjJjMzMiLCJleHAiOjE2MDgyMzcwOTEsInNjb3BlIjoidXNlcnMucHJvZmlsZXM6cmVhZCBhzzRoLnVzZXJfaWRlbnRpdHk6cmVhZCJ9.As1Mz5EbsqkvOekC_mzfsrVEcXX7oPHMzOVzzBQc-zz","token_type":"` + authHeaderBearer + `","expires_in":3599,"refresh_token":"z7696f8a9a92707zbc00a6ab74c474ae9acz8dc4d23125a3z40028a93e65az80","scope":"` + PermissionsIdentity + " " + PermissionsProfile + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPAPIError for mocking requests
type mockHTTPAPIError struct{}

// Do is a mock http request
func (m *mockHTTPAPIError) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"errors":[{"id":"ffb71830-409b-11eb-9032-37efc953c879","status":400,"title":"Bad Request","detail":"Invalid grant: authorization code has expired"}],"jsonapi":{"version":"1.0"}}`)))

	// Default is valid
	return resp, nil
}

// mockHTTPError for mocking requests
type mockHTTPError struct{}

// Do is a mock http request
func (m *mockHTTPError) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusInternalServerError

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))

	// Default is valid
	return resp, nil
}

func TestClient_GetRefreshToken(t *testing.T) {
	t.Parallel()

	t.Run("missing all parameters", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetRefreshToken{})
		assert.NotNil(t, client)
		tokenResponse, err := client.GetRefreshToken(
			context.Background(),
			"",
			"",
			"",
		)
		assert.Error(t, err)
		assert.Nil(t, tokenResponse)
	})

	t.Run("missing client id", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetRefreshToken{})
		assert.NotNil(t, client)
		tokenResponse, err := client.GetRefreshToken(
			context.Background(),
			"",
			"1234567",
			"http://domain.com",
		)
		assert.Error(t, err)
		assert.Nil(t, tokenResponse)
	})

	t.Run("missing auth code", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetRefreshToken{})
		assert.NotNil(t, client)
		tokenResponse, err := client.GetRefreshToken(
			context.Background(),
			"1234567",
			"",
			"http://domain.com",
		)
		assert.Error(t, err)
		assert.Nil(t, tokenResponse)
	})

	t.Run("missing redirect uri", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetRefreshToken{})
		assert.NotNil(t, client)
		tokenResponse, err := client.GetRefreshToken(
			context.Background(),
			"1234567",
			"1234567",
			"",
		)
		assert.Error(t, err)
		assert.Nil(t, tokenResponse)
	})

	t.Run("api error response", func(t *testing.T) {
		client := newTestClient(&mockHTTPAPIError{})
		assert.NotNil(t, client)
		tokenResponse, err := client.GetRefreshToken(
			context.Background(),
			"1234567",
			"1234567",
			"http://domain.com",
		)
		assert.Error(t, err)
		assert.Nil(t, tokenResponse)
	})

	t.Run("http error", func(t *testing.T) {
		client := newTestClient(&mockHTTPError{})
		assert.NotNil(t, client)
		tokenResponse, err := client.GetRefreshToken(
			context.Background(),
			"1234567",
			"1234567",
			"http://domain.com",
		)
		assert.Error(t, err)
		assert.Nil(t, tokenResponse)
	})

	t.Run("valid response", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetRefreshToken{})
		assert.NotNil(t, client)
		tokenResponse, err := client.GetRefreshToken(
			context.Background(),
			"1234567",
			"1234567",
			"http://domain.com",
		)
		assert.NoError(t, err)
		assert.NotNil(t, tokenResponse)
		assert.Equal(t, "zzz0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIzOiIyNjQ1IiwiYXVkIjoiMjhmNzUyZDg0NmRmzjU0YzM3ZWQ1MjBlZTlkYjJjMzMiLCJleHAiOjE2MDgyMzcwOTEsInNjb3BlIjoidXNlcnMucHJvZmlsZXM6cmVhZCBhzzRoLnVzZXJfaWRlbnRpdHk6cmVhZCJ9.As1Mz5EbsqkvOekC_mzfsrVEcXX7oPHMzOVzzBQc-zz", tokenResponse.AccessToken)
		assert.Equal(t, uint32(3599), tokenResponse.ExpiresIn)
		assert.Equal(t, "z7696f8a9a92707zbc00a6ab74c474ae9acz8dc4d23125a3z40028a93e65az80", tokenResponse.RefreshToken)
		assert.Equal(t, PermissionsIdentity+" "+PermissionsProfile, tokenResponse.Scope)
		assert.Equal(t, authHeaderBearer, tokenResponse.TokenType)
	})
}
