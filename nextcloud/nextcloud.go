package nextcloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	cError "github.com/joyread/server/error"
)

// AuthURLRequest struct
type AuthURLRequest struct {
	URL         string
	ClientID    string
	RedirectURI string
}

// GetAuthURL ...
func GetAuthURL(authURLRequest AuthURLRequest) string {
	authURL := fmt.Sprintf("%s/apps/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=write", authURLRequest.URL, authURLRequest.ClientID, authURLRequest.RedirectURI)
	return authURL
}

// AccessTokenRequest struct
type AccessTokenRequest struct {
	URL          string
	ClientID     string
	ClientSecret string
	AuthCode     string
	RedirectURI  string
}

// AccessTokenResponse struct
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// GetAccessToken ...
func GetAccessToken(accessTokenRequest AccessTokenRequest) *AccessTokenResponse {
	body := strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=%s&code=%s&redirect_uri=%s", accessTokenRequest.ClientID, accessTokenRequest.ClientSecret, "authorization_code", accessTokenRequest.AuthCode, accessTokenRequest.RedirectURI))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/apps/oauth2/api/v1/token", accessTokenRequest.URL), body)
	cError.CheckError(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	cError.CheckError(err)

	var accessTokenResponse AccessTokenResponse

	json.NewDecoder(resp.Body).Decode(&accessTokenResponse)
	resp.Body.Close()

	return &accessTokenResponse
}
