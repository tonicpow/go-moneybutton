package moneybutton

import (
	"github.com/tonicpow/go-moneybutton/api"
)

// jsonAPI
type jsonAPI struct {
	Version string `json:"version"`
}

// userIdentityAttributes
type userIdentityAttributes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// userIdentityData
type userIdentityData struct {
	Attributes userIdentityAttributes `json:"attributes"`
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
}

// MoneyButtonUserIdentity is the data returned from MB
type MoneyButtonUserIdentity struct {
	Data    userIdentityData `json:"data"`
	JSONAPI jsonAPI          `json:"jsonapi"`
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
	Attributes userProfileAttributes `json:"attributes"`
	ID         string                `json:"id"`
	Type       string                `json:"type"`
}

// MoneyButtonUserProfile is the user fields returned for the user profile
type MoneyButtonUserProfile struct {
	Data userProfileData `json:"data"`
}

func GetProfile(clientID, oauthToken, redirectURI string) (moneyButtonProfile *MoneyButtonUserIdentity, err error) {
	refreshToken, err := api.GetRefreshToken(clientID, oauthToken, redirectURI)
	if err != nil {
		return nil, err
	}


	// Start the Request
	var request *http.Request
	var jsonValue []byte
	if request, err = http.NewRequestWithContext(context.Background(), http.MethodGet, moneyButtonAPIURL+"auth/user_identity", bytes.NewBuffer(jsonValue)); err != nil {
		nil, err
	}

	// Set refresh token header
	request.Header.Set("Authorization", "Bearer "+refreshTokenResponse.AccessToken)

	// Fire the http Request
	if resp, err = http.DefaultClient.Do(request); err != nil {
		nil, err
	}

	// Close the response body
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		returnError(w, req, a.AppConfig, "failed getting user identity from moneybutton: "+resp.Status+" using refresh token "+refreshTokenResponse.AccessToken+" body "+string(body), "unable to sign in with MoneyButton", http.StatusExpectationFailed)
		return
	}

	// Read the body
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		nil, err
	}

	// Create the response
	userIdentityResponse := new(MoneyButtonUserIdentity)
	if err = json.Unmarshal(body, &userIdentityResponse); err != nil {
		nil, err
	}

	// logger.NoFileData(logger.DEBUG, string(body))

	// Get the user profile using the user Id returned from the profile

	// Start the Request
	if request, err = http.NewRequestWithContext(context.Background(), http.MethodGet, moneyButtonAPIURL+"users/"+userIdentityResponse.Data.ID+"/profile", bytes.NewBuffer(jsonValue)); err != nil {
		
		return nil, err
	}

	// Set refresh token header
	request.Header.Set("Authorization", "Bearer "+refreshTokenResponse.AccessToken)

	// Fire the http Request
	if resp, err = http.DefaultClient.Do(request); err != nil {
		return nil, err
	}

	// Close the response body
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, Error("failed getting user profile from moneybutton: "+resp.Status+" using refresh token "+refreshTokenResponse.AccessToken)
	}

	// Read the body
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}

	moneyButtonProfile := new(MoneyButtonUserProfile)
	if err = json.Unmarshal(body, &moneyButtonProfile); err != nil {
		return nil, err

	} else if moneyButtonProfile == nil {
		return nil, Error("failed to find a MoneyButton user in context")
	} else if len(moneyButtonProfile.Data.Attributes.PrimaryPaymail) == 0 {
		return nil, Error("failed to find a primary paymail for MoneyButton user")
	}

}
