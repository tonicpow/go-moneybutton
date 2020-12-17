package moneybutton

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
{
  "data": {
    "type": "profiles",
    "id": "123",
    "attributes": {
      "created-at": "2019-03-26T17:33:42.788Z",
      "name": "MrZ",
      "default-currency": "USD",
      "default-language": "en",
      "bio": "I like Money Button.",
      "primary-paymail": "mrz@moneybutton.com",
      "avatar-url": "https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon"
    }
  }
}
*/

// GetProfile returns profile info for the specified user
//
// Specs: https://docs.moneybutton.com/docs/api-rest-user-profile.html
func (c *Client) GetProfile(ctx context.Context, userID, accessToken string) (*UserProfile, error) {

	// Check required parameters
	if len(accessToken) == 0 {
		return nil, fmt.Errorf("missing required parameter: %s", "accessToken")
	} else if len(userID) == 0 {
		return nil, fmt.Errorf("missing required parameter: %s", "userID")
	}

	// Fire the request
	response := httpRequest(
		ctx,
		c,
		&httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			Token:          accessToken,
			URL:            fmt.Sprintf(endpointUserProfile, userID),
		},
	)

	// Error in request?
	if response.Error != nil {
		return nil, response.Error
	}

	// Create the response
	profile := new(UserProfile)
	if err := json.Unmarshal(response.BodyContents, &profile); err != nil {
		return nil, err
	}
	return profile, nil
}
