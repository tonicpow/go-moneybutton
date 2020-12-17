package moneybutton

// RefreshTokenResponse is used to get a refresh token for getting
// user information from the moneybutton API
//
// Specs: https://docs.moneybutton.com/docs/api-oauth-endpoints.html
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"` // "Bearer"
	ExpiresIn    uint32 `json:"expires_in"` // 3600  (default is 1 hour)
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// UserIdentity is the user identity
//
// Specs: https://docs.moneybutton.com/docs/api-rest-user-identity.html
type UserIdentity struct {
	Data    *userIdentityData `json:"data"`
	JSONAPI *jsonAPIVersion   `json:"jsonapi"`
}

// userIdentityAttributes is used in identity data
type userIdentityAttributes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// userIdentityData is the identity data
type userIdentityData struct {
	Attributes *userIdentityAttributes `json:"attributes"`
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
}

// UserProfile is the user fields returned for the user profile
//
// Specs: https://docs.moneybutton.com/docs/api-rest-user-profile.html
type UserProfile struct {
	Data *userProfileData `json:"data"`
}

// userProfileAttributes
type userProfileAttributes struct {
	AvatarURL       string `json:"avatar-url"`
	Bio             string `json:"bio"`
	CreatedAt       string `json:"created-at"`
	DefaultCurrency string `json:"default-currency"`
	DefaultLanguage string `json:"default-language"`
	Name            string `json:"name"`
	PrimaryPaymail  string `json:"primary-paymail"`
}

// userProfileData
type userProfileData struct {
	Attributes *userProfileAttributes `json:"attributes"`
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
}

/*
{
  "errors": [
    {
      "id": "ffb71830-409b-11eb-9032-37efc953c879",
      "status": 400,
      "title": "Bad Request",
      "detail": "Invalid grant: authorization code has expired"
    }
  ],
  "jsonapi": {
    "version": "1.0"
  }
}
*/

// errorResponse is the error response from the MoneyButton API
type errorResponse struct {
	Errors  []*apiError     `json:"errors"`
	JSONAPI *jsonAPIVersion `json:"jsonapi"`
}

// apiError is the individual error
type apiError struct {
	Detail string `json:"detail"`
	ID     string `json:"id"`
	Status int    `json:"status"`
	Title  string `json:"title"`
}

// jsonAPIVersion is the API version
type jsonAPIVersion struct {
	Version string `json:"version"`
}
