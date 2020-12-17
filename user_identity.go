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
    "id": "123",
    "type": "user_identities",
    "attributes": {
      "id": "123",
      "name": "MrZ"
    }
  },
  "jsonapi": {
    "version": "1.0"
  }
}
*/

// GetUserIdentity returns the minimum data to identify a user
//
// Specs: https://docs.moneybutton.com/docs/api-rest-user-identity.html
func (c *Client) GetUserIdentity(ctx context.Context, accessToken string) (*UserIdentity, error) {

	// Check required parameters
	if len(accessToken) == 0 {
		return nil, fmt.Errorf("missing required parameter: %s", "accessToken")
	}

	// Fire the request
	response := httpRequest(
		ctx,
		c,
		&httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			Token:          accessToken,
			URL:            endpointUserIdentity,
		},
	)

	// Error in request?
	if response.Error != nil {
		return nil, response.Error
	}

	// Create the response
	identity := new(UserIdentity)
	if err := json.Unmarshal(response.BodyContents, &identity); err != nil {
		return nil, err
	}
	return identity, nil
}
