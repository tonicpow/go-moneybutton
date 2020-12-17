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

// mockHTTPGetUserProfile for mocking requests
type mockHTTPGetUserProfile struct{}

// Do is a mock http request
func (m *mockHTTPGetUserProfile) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	if req.URL.String() == fmt.Sprintf(endpointUserProfile, "123") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"data":{"type":"profiles","id":"123","attributes":{"created-at":"2019-03-26T17:33:42.788Z","name":"MrZ","default-currency":"USD","default-language":"en","bio":"I like Money Button.","primary-paymail":"mrz@moneybutton.com","avatar-url":"https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon"}}}`)))
	}

	// Default is valid
	return resp, nil
}

func TestClient_GetProfile(t *testing.T) {
	t.Parallel()

	t.Run("missing all parameters", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetUserProfile{})
		assert.NotNil(t, client)
		profile, err := client.GetProfile(
			context.Background(),
			"",
			"",
		)
		assert.Error(t, err)
		assert.Nil(t, profile)
	})

	t.Run("missing user id", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetUserProfile{})
		assert.NotNil(t, client)
		profile, err := client.GetProfile(
			context.Background(),
			"",
			"1234567",
		)
		assert.Error(t, err)
		assert.Nil(t, profile)
	})

	t.Run("missing access token", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetUserProfile{})
		assert.NotNil(t, client)
		profile, err := client.GetProfile(
			context.Background(),
			"123",
			"",
		)
		assert.Error(t, err)
		assert.Nil(t, profile)
	})

	t.Run("api error response", func(t *testing.T) {
		client := newTestClient(&mockHTTPAPIError{})
		assert.NotNil(t, client)
		profile, err := client.GetProfile(
			context.Background(),
			"123",
			"1234567",
		)
		assert.Error(t, err)
		assert.Nil(t, profile)
	})

	t.Run("http error", func(t *testing.T) {
		client := newTestClient(&mockHTTPError{})
		assert.NotNil(t, client)
		profile, err := client.GetProfile(
			context.Background(),
			"123",
			"1234567",
		)
		assert.Error(t, err)
		assert.Nil(t, profile)
	})

	t.Run("valid response", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetUserProfile{})
		assert.NotNil(t, client)
		profile, err := client.GetProfile(
			context.Background(),
			"123",
			"1234567",
		)
		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, "profiles", profile.Data.Type)
		assert.Equal(t, "123", profile.Data.ID)
		assert.Equal(t, "MrZ", profile.Data.Attributes.Name)
		assert.Equal(t, "https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon", profile.Data.Attributes.AvatarURL)
		assert.Equal(t, "I like Money Button.", profile.Data.Attributes.Bio)
		assert.Equal(t, "2019-03-26T17:33:42.788Z", profile.Data.Attributes.CreatedAt)
		assert.Equal(t, "USD", profile.Data.Attributes.DefaultCurrency)
		assert.Equal(t, "en", profile.Data.Attributes.DefaultLanguage)
		assert.Equal(t, "mrz@moneybutton.com", profile.Data.Attributes.PrimaryPaymail)
	})
}
