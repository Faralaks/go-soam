package vk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	. "tools"
)

func GetOauthLink() string {
	url := fmt.Sprintf("https://oauth.vk.com/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		Config.OauthClientID, Config.OauthRedirectURL, "account", "skakoystate")
	return url
}

func GetUserOauthToken(code string) string {
	url := fmt.Sprintf("https://oauth.vk.com/access_token?grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s",
		code, Config.OauthRedirectURL, Config.OauthClientID, Config.OauthKey)
	VPrint(url)
	req, _ := http.NewRequest("POST", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err.Error())
		return ""
	}
	defer resp.Body.Close()

	token := struct {
		AccessToken string `json:"access_token"`
	}{}
	bytes, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(bytes, &token)
	VPrint(string(bytes))
	url = fmt.Sprintf("https://api.vk.com/method/%s?v=5.124&access_token=%s", "users.get", token.AccessToken)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
		return ""
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err.Error())
		return ""
	}
	defer resp.Body.Close()
	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
		return ""
	}
	VPrint(string(bytes))
	return string(bytes)
}
