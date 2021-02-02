package moneybutton

const (

	// version is the current package version
	version = "v0.0.2"

	// defaultUserAgent is the default user agent for all requests
	defaultUserAgent string = "go-moneybutton: " + version

	// apiVersion of the MoneyButton API
	apiVersion = "v1"

	// grants for oAuth
	grantTypeAuthorizationCode  = "authorization_code"
	grantTypeRefreshAccessToken = "refresh_token"

	// endpoints
	endpointToken        = OauthURL + "token"
	endpointUserIdentity = APIURL + "auth/user_identity"
	endpointUserProfile  = APIURL + "users/%s/profile" // requires fmt.Sprintf(endpointUserProfile,userID)

	// authorization header
	authHeaderBearer = "Bearer"
)

// Public constants used for MoneyButton
//
// Specs: https://docs.moneybutton.com/docs/api-overview.html
const (
	APIURL   = "https://www.moneybutton.com/api/" + apiVersion + "/"
	OauthURL = "https://www.moneybutton.com/oauth/" + apiVersion + "/"

	// MoneyButton oAuth Permissions
	PermissionsBalance  = "users.balance:read"
	PermissionsIdentity = "auth.user_identity:read"
	PermissionsProfile  = "users.profiles:read"
)
