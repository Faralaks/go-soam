package handlers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"time"
	. "tools"
)

var SignUp = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	ege, err := strconv.Atoi(TrimStr(r.FormValue("ege"), 2))
	grade, err := strconv.Atoi(TrimStr(r.FormValue("grade"), 2))
	name := TrimStr(r.FormValue("username"), 30)

	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Недопустимое значение в одном из полей: Возраст, Класс | " + err.Error()}.Send(w)
		return
	}

	owner := r.Header.Get("owner")
	userUid, err := p.ObjectIDFromHex(owner)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось преобразовать Uid пользователя в ObjectID | " + err.Error()}.Send(w)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, err = UsersCol.UpdateOne(ctx, bson.M{"_id": userUid, "status": TesteeStatus}, bson.D{{"$set", bson.M{
		"step":      0,
		"tests":     Config.TestList,
		"birthYear": uint16(time.Now().Year() - ege),
		"ege":       uint8(ege),
		"grade":     uint8(grade),
		"name":      NewB64String(name),
	}}})

	if err != nil {
		errCode := err.(mongo.WriteException).WriteErrors[0].Code
		if errCode == 11000 {
			JsonMsg{Kind: DuplicateKeyKind, Field: "newLogin"}.Send(w)
		} else {
			JsonMsg{Kind: FatalKind, Msg: "Не удалось сохранить нового пользователя в базу данных | " + err.Error()}.Send(w)
		}
		return
	}

	JsonMsg{Kind: SucKind}.Send(w)
})
