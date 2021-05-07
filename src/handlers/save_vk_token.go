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
		VPrint("code пуст!")
	}

	url := fmt.Sprintf("https://oauth.vk.com/access_token?grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s",
		code, Config.OauthRedirectURL, Config.OauthClientID, Config.OauthKey)

	req, _ := http.NewRequest("POST", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		VPrint(err.Error())
		return
	}
	defer resp.Body.Close()

	var token VKTokenData
	bytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bytes, &token)
	if err != nil {
		VPrint("Проблема при анмаршалинге токена | " + err.Error() + " | body: " + string(bytes))
		return
	}

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
