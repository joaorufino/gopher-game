package linkedin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	clientID     = "77xiyfrl47693v"
	clientSecret = "WPL_AP1.ogdI1noPbsdWGSrH.PZEG8g=="
	redirectURI  = "http://localhost:8081/callback"
	authURL      = "https://www.linkedin.com/oauth/v2/authorization"
	tokenURL     = "https://www.linkedin.com/oauth/v2/accessToken"
	profileURL   = "https://api.linkedin.com/v2/me"
)

// GetAuthURL generates the URL for user authorization.
func GetAuthURL(state string) string {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("state", state)
	params.Add("scope", "r_liteprofile")

	return fmt.Sprintf("%s?%s", authURL, params.Encode())
}

// AuthToken represents the structure of the OAuth token response.
type AuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// Profile represents the structure of the LinkedIn profile response.
type Profile struct {
	FirstName string `json:"localizedFirstName"`
	LastName  string `json:"localizedLastName"`
	Headline  string `json:"headline"`
}

// GetAccessToken exchanges the authorization code for an access token.
func GetAccessToken(authCode string) (*AuthToken, error) {
	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("code", authCode)
	params.Add("redirect_uri", redirectURI)
	params.Add("client_id", clientID)
	params.Add("client_secret", clientSecret)

	resp, err := http.PostForm(tokenURL, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token AuthToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

// GetProfile fetches the LinkedIn profile using the access token.
func GetProfile(accessToken string) (*Profile, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", profileURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}
	return &profile, nil
}
