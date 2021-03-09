package hendlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
	. "tools"
)

type TesteeListResponse struct {
	TesteeList []*Testee `json:"testeeList"`
}

var Get_testee_list = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var response TesteeListResponse
	owner := r.Header.Get("owner")
	if r.Header.Get("status") == AdminStatus {
		owner = TrimStr(r.FormValue("psyUid"), 30)
	}
	if owner == "" {
		JsonMsg{Kind: FatalKind, Msg: "Не был получен uid  психолога"}.SendMsg(w)
		return
	}
	grade := NewB64UpString(TrimStr(r.FormValue("grade"), 30))
	if grade == "" {
		JsonMsg{Kind: FatalKind, Msg: "Не был получен класс испытуемых | "}.SendMsg(w)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 4*time.Second)
	cur, err := UsersCol.Find(ctx, bson.M{"status": TesteeStatus, "owner": owner, "grade": grade})
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось извлечь список испытуемых | " + err.Error()}.SendMsg(w)
		return
	}
	if cur == nil {
		JsonMsg{Kind: FatalKind, Msg: "ТАКОГО НЕ МОЖЕТ СЛУЧИТСЯ, ПОТОМУ ЧТО НЕ МОЖЕТ. КУРСОР, МАТЬ ЕГО, СТАЛ NIL"}.SendMsg(w)
		return
	}
	defer cur.Close(context.TODO())

	ctx, _ = context.WithTimeout(context.Background(), 6*time.Second)
	for cur.Next(ctx) {
		var elem Testee
		err := cur.Decode(&elem)
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: "В процессе декодирования испытуемых произошла ошибка | " + err.Error()}.SendMsg(w)
			return
		}
		elem.Pas, err = Decrypt(elem.Pas)
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: fmt.Sprintf("В процессе дешифрования пароля пользователя %v произошла ошибка | %v", elem.Login.Decode(), err.Error())}.SendMsg(w)
			return
		}
		response.TesteeList = append(response.TesteeList, &elem)
	}

	if err := cur.Err(); err != nil {
		JsonMsg{Kind: FatalKind, Msg: "У курсора произошла ошибка | " + err.Error()}.SendMsg(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "В процессе конвертации в json возникла ошибка | " + err.Error()}.SendMsg(w)
		return
	}
})
