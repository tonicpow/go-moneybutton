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

// mockHTTPGetUserIdentity for mocking requests
type mockHTTPGetUserIdentity struct{}

// Do is a mock http request
func (m *mockHTTPGetUserIdentity) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	if req.URL.String() == endpointUserIdentity {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"data":{"id":"123","type":"user_identities","attributes":{"id":"123","name":"MrZ"}},"jsonapi":{"version":"1.0"}}`)))
	}

	// Default is valid
	return resp, nil
}

func TestClient_GetUserIdentity(t *testing.T) {
	t.Parallel()

	t.Run("missing access token", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetUserIdentity{})
		assert.NotNil(t, client)
		identity, err := client.GetUserIdentity(
			context.Background(),
			"",
		)
		assert.Error(t, err)
		assert.Nil(t, identity)
	})

	t.Run("api error response", func(t *testing.T) {
		client := newTestClient(&mockHTTPAPIError{})
		assert.NotNil(t, client)
		identity, err := client.GetUserIdentity(
			context.Background(),
			"1234567",
		)
		assert.Error(t, err)
		assert.Nil(t, identity)
	})

	t.Run("http error", func(t *testing.T) {
		client := newTestClient(&mockHTTPError{})
		assert.NotNil(t, client)
		identity, err := client.GetUserIdentity(
			context.Background(),
			"1234567",
		)
		assert.Error(t, err)
		assert.Nil(t, identity)
	})

	t.Run("valid response", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetUserIdentity{})
		assert.NotNil(t, client)
		identity, err := client.GetUserIdentity(
			context.Background(),
			"1234567",
		)
		assert.NoError(t, err)
		assert.NotNil(t, identity)
		assert.Equal(t, "1.0", identity.JSONAPI.Version)
		assert.Equal(t, "user_identities", identity.Data.Type)
		assert.Equal(t, "123", identity.Data.ID)
		assert.Equal(t, "123", identity.Data.Attributes.ID)
		assert.Equal(t, "MrZ", identity.Data.Attributes.Name)
	})
}
