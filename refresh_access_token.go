package moneybutton

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
{
  "access_token": "eyJ0eXAiOi...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "scope": "users.profiles:read auth.user_identity:read",
  "refresh_token": "379fc72ba..."
}
*/

// RefreshAccessToken will refresh an existing access token
//
// Specs: https://docs.moneybutton.com/docs/api-oauth-endpoints.html#requesting-the-refresh-token
func (c *Client) RefreshAccessToken(ctx context.Context, clientID, accessToken string) (*RefreshTokenResponse, error) {

	// Check required parameters
	if len(clientID) == 0 {
		return nil, fmt.Errorf("missing required parameter: %s", "clientID")
	} else if len(accessToken) == 0 {
		return nil, fmt.Errorf("missing required parameter: %s", "accessToken")
	}

	// Fire the request
	response := httpRequest(
		ctx,
		c,
		&httpPayload{
			Data: `grant_type=` + grantTypeRefreshAccessToken + `&client_id=` + clientID +
				`&refresh_token=` + accessToken,
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodPost,
			URL:            endpointToken,
		},
	)

	// Error in request?
	if response.Error != nil {
		return nil, response.Error
	}

	// Create the response
	refreshTokenResponse := new(RefreshTokenResponse)
	if err := json.Unmarshal(response.BodyContents, &refreshTokenResponse); err != nil {
		return nil, err
	}
	return refreshTokenResponse, nil
}
