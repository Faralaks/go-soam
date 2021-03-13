package handlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	. "tools"
)

type UserDataResponse struct {
	UserData *FullUser `json:"userData"`
}

var Get_user_data = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var response UserDataResponse
	curUserUid, err := p.ObjectIDFromHex(r.Header.Get("owner"))

	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось преобразовать uid в ObjectID | " + err.Error()}.Send(w)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = UsersCol.FindOne(ctx, bson.M{"_id": curUserUid}).Decode(&response.UserData)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось получить данные пользователя | " + err.Error()}.Send(w)
		return
	}

	response.UserData.Pas = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "В процессе конвертации в json возникла ошибка | " + err.Error()}.Send(w)
		return
	}

})
