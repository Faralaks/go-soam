package hendlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	. "tools"
)

type UserDataResponse struct {
	UserData *MultiUser `json:"userData"`
}

var Get_user_data = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var response UserDataResponse
	curUserUid, err := p.ObjectIDFromHex(r.Header.Get("owner"))
	if r.Header.Get("status") == AdminStatus && TrimStr(r.FormValue("isMy"), 4) != "true" {
		curUserUid, err = p.ObjectIDFromHex(TrimStr(r.FormValue("psyUid"), 30))
	}
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось преобразовать uid в ObjectID | " + err.Error()}.SendMsg(w)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = UsersCol.FindOne(ctx, bson.M{"_id": curUserUid}).Decode(&response.UserData)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось получить данные пользователя | " + err.Error()}.SendMsg(w)
		return
	}
	//fmt.Printf("\n%v\n\n", response.UserData.Available)
	if r.Header.Get("status") == AdminStatus && response.UserData.Status != AdminStatus {
		response.UserData.Pas, err = Decrypt(response.UserData.Pas)
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: fmt.Sprintf("Не удалось расшифровать пароль пользователя %v | "+err.Error(), response.UserData.Login.Decode())}.SendMsg(w)
			return
		}
	} else {
		response.UserData.Pas = ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "В процессе конвертации в json возникла ошибка | " + err.Error()}.SendMsg(w)
		return
	}

})
