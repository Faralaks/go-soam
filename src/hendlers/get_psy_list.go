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

type PsyListResponse struct {
	PsyList []*Psy `json:"psyList"`
}

var Get_psy_list = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 4*time.Second)
	cur, err := UsersCol.Find(ctx, bson.M{"status": PsyStatus})
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось извлечь список психологов | " + err.Error()}.SendMsg(w)
		return
	}
	if cur == nil {
		JsonMsg{Kind: FatalKind, Msg: "ТАКОГО НЕ МОЖЕТ СЛУЧИТСЯ, ПОТОМУ ЧТО НЕ МОЖЕТ. КУРСОР, МАТЬ ЕГО, СТАЛ NIL"}.SendMsg(w)
		return
	}
	defer cur.Close(context.TODO())

	var response PsyListResponse
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	for cur.Next(ctx) {
		var elem Psy
		err := cur.Decode(&elem)
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: "В процессе декодирования психологов произошла ошибка | " + err.Error()}.SendMsg(w)
			return
		}
		elem.Pas, err = Decrypt(elem.Pas)
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: fmt.Sprintf("В процессе дешифрования пароля пользователя %v произошла ошибка | %v", elem.Login.Decode(), err.Error())}.SendMsg(w)
			return
		}
		response.PsyList = append(response.PsyList, &elem)
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
