package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	. "tools"
)

var SaveVKToken = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	code := TrimStr(r.URL.Query()["code"][0], 30)
	if code == "" {
		JsonMsg{Kind: FatalKind, Msg: "Не был передан code"}.Send(w)
		return
	}

	url := fmt.Sprintf("https://oauth.vk.com/access_token?grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s",
		code, Config.OauthRedirectURL, Config.OauthClientID, Config.OauthKey)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось создать запрос к ВК для получения токена | " + err.Error()}.Send(w)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось соверишть запрос к ВК для получения токена | " + err.Error()}.Send(w)
		return
	}
	defer resp.Body.Close()

	var token VKTokenData
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Проблема при чтении пакета с токеном | " + err.Error()}.Send(w)
	}

	err = json.Unmarshal(bytes, &token)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Проблема при анмаршалинге токена | " + err.Error()}.Send(w)
		return
	}

	url = fmt.Sprintf("https://api.vk.com/method/users.get?v=5.124&access_token=%s", token.Token)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось создать запрос к ВК для получения личных данных | " + err.Error()}.Send(w)
		return
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось совершить запрос к ВК для получения личных данных | " + err.Error()}.Send(w)
		return
	}
	defer resp.Body.Close()

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Проблема при чтении пакета с данными пользователя | " + err.Error()}.Send(w)
	}

	var userData = struct {
		Response []struct {
			Name string `json:"first_name"`
		} `json:"response"`
	}{}
	err = json.Unmarshal(bytes, &userData)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Проблема при анмаршалинге пользовательских данных | " + err.Error()}.Send(w)
		VPrint(string(bytes))
		return
	}

	token.UserName = userData.Response[0].Name

	newUser := User{
		Uid:         p.NewObjectID(),
		Login:       NewB64LowString(strconv.Itoa(token.UserId)),
		Status:      TesteeStatus,
		CreatedDate: CurUtcStamp(),
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	_, err = UsersCol.InsertOne(ctx, newUser)
	if err != nil {
		errCode := err.(mongo.WriteException).WriteErrors[0].Code
		if errCode == 11000 {
			JsonMsg{Kind: DuplicateKeyKind, Field: "newLogin"}.Send(w)
		} else {
			JsonMsg{Kind: FatalKind, Msg: "Не удалось сохранить нового пользователя в базу данных | " + err.Error()}.Send(w)
		}
		return
	}

	signedAt, signedRt, err := CreateTokens(newUser.Uid.Hex(), newUser.Status)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось создать токены | " + err.Error()}.Send(w)
		return
	}
	SetLoginCookies(w, signedAt, signedRt)

	w.Header().Set("Content-Type", "text/html")
	main := Config.CurPath + "/templates/signup.html"
	base := Config.CurPath + "/templates/base.html"

	tmpl, err := template.ParseFiles(main, base)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = tmpl.ExecuteTemplate(w, "signup", token)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

})
