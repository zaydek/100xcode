package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/oauth1"
)

////////////////////////////////////////////////////////////////////////////////

type OAuth1AuthenticationParameters struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessKey      string
	AccessSecret   string
}

type OAuth1Authentication struct {
	*http.Client
}

func newOAuth1Authentication(params OAuth1AuthenticationParameters) *OAuth1Authentication {
	config := oauth1.NewConfig(params.ConsumerKey, params.ConsumerSecret)
	httpClient := config.Client(oauth1.NoContext, oauth1.NewToken(params.AccessKey, params.AccessSecret))
	if httpClient == nil {
		return nil
	}
	return &OAuth1Authentication{httpClient}
}

func (a *OAuth1Authentication) Get(url string) ([]byte, error) {
	resp, err := a.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

////////////////////////////////////////////////////////////////////////////////

type TwitterOAuth1Authentication struct {
	*OAuth1Authentication
}

func newTwitterOAuth1Authentication(params OAuth1AuthenticationParameters) *TwitterOAuth1Authentication {
	oauth1Auth := newOAuth1Authentication(OAuth1AuthenticationParameters{
		ConsumerKey:    CONSUMER_KEY,
		ConsumerSecret: CONSUMER_SECRET,
		AccessKey:      ACCESS_KEY,
		AccessSecret:   ACCESS_SECRET,
	})
	if oauth1Auth == nil {
		return nil
	}
	return &TwitterOAuth1Authentication{oauth1Auth}
}

func (a *TwitterOAuth1Authentication) GetBlockedScreenNames() ([]string, error) {
	// https://developer.twitter.com/en/docs/twitter-api/v1/accounts-and-users/mute-block-report-users/api-reference/get-blocks-ids
	data, err := a.OAuth1Authentication.Get("https://api.twitter.com/1.1/blocks/list.json?include_entities=false&skip_status=true")
	if err != nil {
		return nil, err
	}
	// https://mholt.github.io/json-to-go/
	var blockedResponse struct {
		Users []struct {
			ScreenName string `json:"screen_name"`
		} `json:"users"`
	}
	if err := json.Unmarshal(data, &blockedResponse); err != nil {
		return nil, err
	}
	var screenNames []string
	for _, user := range blockedResponse.Users {
		screenNames = append(screenNames, user.ScreenName)
	}
	return screenNames, nil
}
