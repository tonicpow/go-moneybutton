package moneybutton

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
{
  "access_token": "eyJ0eXAiOiJKV1....",
  "token_type": "Bearer",
  "expires_in": 3599,
  "refresh_token": "e7cb3bac6....",
  "scope": "users.profiles:read auth.user_identity:read"
}
*/

// GetRefreshToken will get a new refresh token given an auth code
//
// Specs: https://docs.moneybutton.com/docs/api-oauth-endpoints.html#requesting-the-refresh-token
func (c *Client) GetRefreshToken(ctx context.Context, clientID, authCode,
	redirectURI string) (*RefreshTokenResponse, error) {

	// Check required parameters
	if len(clientID) == 0 {
		return nil, fmt.Errorf("missing required parameter: %s", "clientID")
	} else if len(authCode) == 0 {
		return nil, fmt.Errorf("missing required parameter: %s", "authCode")
	} else if len(redirectURI) == 0 {
		return nil, fmt.Errorf("missing required parameter: %s", "redirectURI")
	}

	// Fire the request
	response := httpRequest(
		ctx,
		c,
		&httpPayload{
			Data: `grant_type=` + grantTypeAuthorizationCode + `&client_id=` + clientID +
				`&code=` + authCode + `&redirect_uri=` + redirectURI,
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
