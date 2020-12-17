package moneybutton

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// RequestResponse is the response from a request
type RequestResponse struct {
	BodyContents []byte `json:"body_contents"` // Raw body response
	Error        error  `json:"error"`         // If an error occurs
	Method       string `json:"method"`        // Method is the HTTP method used
	PostData     string `json:"post_data"`     // PostData is the post data submitted if POST/PUT request
	StatusCode   int    `json:"status_code"`   // StatusCode is the last code from the request
	URL          string `json:"url"`           // URL is used for the request
}

// httpPayload is used for a httpRequest
type httpPayload struct {
	Data           string `json:"data"`
	ExpectedStatus int    `json:"expected_status"`
	Method         string `json:"method"`
	Token          string `json:"token"`
	URL            string `json:"url"`
}

// httpRequest is a generic request wrapper that can be used without constraints
func httpRequest(ctx context.Context, client *Client,
	payload *httpPayload) (response *RequestResponse) {

	// Set reader
	var bodyReader io.Reader

	// Start the response
	response = new(RequestResponse)

	// Add post data if applicable
	if payload.Method == http.MethodPost || payload.Method == http.MethodPut {
		bodyReader = strings.NewReader(payload.Data)
		response.PostData = payload.Data
	}

	// Store for debugging purposes
	response.Method = payload.Method
	response.URL = payload.URL

	// Start the request
	var request *http.Request
	if request, response.Error = http.NewRequestWithContext(
		ctx, payload.Method, payload.URL, bodyReader,
	); response.Error != nil {
		return
	}

	// Change the header (user agent is in case they block default Go user agents)
	request.Header.Set("User-Agent", client.Options.UserAgent)

	// Set the content type on Method
	if payload.Method == http.MethodPost || payload.Method == http.MethodPut {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// Set the token if found
	if len(payload.Token) > 0 {
		request.Header.Set("Authorization", authHeaderBearer+" "+payload.Token)
	}

	// Fire the http request
	var resp *http.Response
	if resp, response.Error = client.httpClient.Do(request); response.Error != nil {
		if resp != nil {
			response.StatusCode = resp.StatusCode
		}
		return
	}

	// Close the response body
	defer func() {
		_ = resp.Body.Close()
	}()

	// Set the status
	response.StatusCode = resp.StatusCode

	// Read the body
	if response.BodyContents, response.Error = ioutil.ReadAll(resp.Body); response.Error != nil {
		return
	}

	// Status does not match as expected
	if resp.StatusCode != payload.ExpectedStatus {

		// Set the error message (return 1 for now)
		if len(response.BodyContents) > 0 {
			errorMsg := new(errorResponse)
			if response.Error = json.Unmarshal(
				response.BodyContents, &errorMsg,
			); response.Error != nil {
				return
			}
			// todo: this needs some love (supporting multiple errors)
			errString := ""
			for _, err := range errorMsg.Errors {
				errString += "error: " + err.Detail + " "
			}

			response.Error = fmt.Errorf("%s", errString)
			return
		}

		// No error message found, set default error message
		response.Error = fmt.Errorf("request failed with status code: %d", resp.StatusCode)
		return
	}

	return
}
