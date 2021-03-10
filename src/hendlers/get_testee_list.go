package hendlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
	. "tools"
)

type TesteeListResponse struct {
	TesteeList []*FullUser `json:"testeeList"`
}

var Get_testee_list = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var response TesteeListResponse

	ctx, _ := context.WithTimeout(context.Background(), 4*time.Second)
	cur, err := UsersCol.Find(ctx, bson.M{"status": TesteeStatus})
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось извлечь список испытуемых | " + err.Error()}.SendMsg(w)
		return
	}
	defer cur.Close(context.TODO())

	ctx, _ = context.WithTimeout(context.Background(), 6*time.Second)
	for cur.Next(ctx) {
		var elem FullUser
		err := cur.Decode(&elem)
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: "В процессе декодирования испытуемых произошла ошибка | " + err.Error()}.SendMsg(w)
			return
		}
		elem.Pas = elem.Pas
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
