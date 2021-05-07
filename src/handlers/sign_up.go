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
		"Step":      0,
		"Tests":     Config.TestList,
		"BirthYear": uint16(time.Now().Year() - ege),
		"Ege":       uint8(ege),
		"Grade":     uint8(grade),
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
