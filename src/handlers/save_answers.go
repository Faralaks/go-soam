package handlers

import (
	. "blank_handlers"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"
	. "tools"
)

var Save_answers = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	curUserUid, err := p.ObjectIDFromHex(r.Header.Get("owner"))
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось преобразовать uid в ObjectID | " + err.Error()}.Send(w)
		return
	}

	ansList := strings.Split(r.FormValue("answers"), "&")

	var user FullUser
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = UsersCol.FindOne(ctx, bson.M{"_id": curUserUid}).Decode(&user)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось получить данные пользователя | " + err.Error()}.Send(w)
		return
	}

	step := Config.TestList[user.Step]

	err = BlankHandlers[step](ansList, curUserUid)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось посчитать и сохранить результаты | " + err.Error()}.Send(w)
		return
	}

	JsonMsg{Kind: SucKind}.Send(w)

})
