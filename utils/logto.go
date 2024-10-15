package utils

import (
	"encoding/json"
	"github.com/logto-io/go/client"
	"github.com/logto-io/go/core"
	"io"
	"net/http"
	"net/url"
)

type UserInfos struct {
	Sub               string     `json:"sub"`
	Name              string     `json:"name"`
	FamilyName        string     `json:"family_name"`
	GivenName         string     `json:"given_name"`
	MiddleName        string     `json:"middle_name"`
	Nickname          string     `json:"nickname"`
	PreferredUsername string     `json:"preferred_username"`
	Profile           string     `json:"profile"`
	Website           string     `json:"website"`
	Picture           string     `json:"picture"`
	Gender            string     `json:"gender"`
	Birthdate         string     `json:"birthdate"`
	Zoneinfo          string     `json:"zoneinfo"`
	Locale            string     `json:"locale"`
	UpdatedAt         int64      `json:"updated_at"`
	Username          string     `json:"username"`
	CreatedAt         int64      `json:"created_at"`
	Email             string     `json:"email"`
	EmailVerified     bool       `json:"email_verified"`
	CustomData        CustomData `json:"custom_data"`
}

type CustomData struct {
	CanInvite bool `json:"can_invite"`
}

func FetchUserInfos(logtoClient *client.LogtoClient, logtoConfig *client.LogtoConfig) (userProfile UserInfos, err error) {
	client := &http.Client{}

	discoveryEndpoint, err := url.JoinPath(logtoConfig.Endpoint, "/oidc/.well-known/openid-configuration")
	if err != nil {
		return UserInfos{}, err
	}
	oidcConfig, err := core.FetchOidcConfig(client, discoveryEndpoint)
	if err != nil {
		return UserInfos{}, err
	}

	accessToken, err := logtoClient.GetAccessToken("")
	if err != nil {
		return UserInfos{}, err
	}

	request, err := http.NewRequest("GET", oidcConfig.UserinfoEndpoint, nil)
	if err != nil {
		return UserInfos{}, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken.Token)

	response, err := client.Do(request)
	if err != nil {
		return UserInfos{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return UserInfos{}, err
	}

	var userInfos UserInfos
	err = json.Unmarshal(body, &userInfos)
	if err != nil {
		return UserInfos{}, err
	}
	return userInfos, nil
}
