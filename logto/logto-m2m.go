package logto

import (
	"LogtoUserProfile/globals"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

type ProfilePayload struct {
	GivenName  string `json:"givenName,omitempty"`
	FamilyName string `json:"familyName,omitempty"`
}

type PatchProfilePayload struct {
	Profile ProfilePayload `json:"profile,omitempty"`
	Avatar  string         `json:"avatar,omitempty"`
	Name    string         `json:"name,omitempty"`
}

type PatchProfilePayloadLite struct {
	Avatar string `json:"avatar,omitempty"`
	Name   string `json:"name,omitempty"`
}

var accessToken struct {
	Token      string
	ExpiryDate int64
}

func getAccessToken() (string, error) {
	actualTime := time.Now().Unix()
	if accessToken.Token != "" && actualTime < accessToken.ExpiryDate {
		return accessToken.Token, nil
	}

	apiUrl, err := url.JoinPath(globals.Configuration.Logto.Endpoint, "/oidc/token")
	if err != nil {
		return "", err
	}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("resource", "https://default.logto.app/api")
	data.Set("scope", "all")

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(globals.Configuration.Logto.M2MAppId, globals.Configuration.Logto.M2MAppSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch access token")
	}

	var responseJSON struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	err = json.NewDecoder(resp.Body).Decode(&responseJSON)
	if err != nil {
		return "", err
	}

	actualTime = time.Now().Unix()
	accessToken.Token = responseJSON.AccessToken
	accessToken.ExpiryDate = actualTime + responseJSON.ExpiresIn

	return accessToken.Token, nil
}

func PatchUserProfile(sub string, payload interface{}) error {
	token, err := getAccessToken()
	if err != nil {
		return err
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	endpoint, err := url.JoinPath(globals.Configuration.Logto.Endpoint, "/api/users/", sub)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Vérifier la réponse
	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		return errors.New(resp.Status)
	}
}
