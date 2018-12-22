/*
	Copyright (C) 2018 Nirmal Almara

    This file is part of Joyread.

    Joyread is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    Joyread is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
	along with Joyread.  If not, see <https://www.gnu.org/licenses/>.
*/

package nextcloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	cError "gitlab.com/joyread/ultimate/error"
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
