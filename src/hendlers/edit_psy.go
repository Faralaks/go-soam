package hendlers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"
	. "tools"
)

var Edit_psy = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	psyUid, err := p.ObjectIDFromHex(TrimStr(r.FormValue("psyUid"), 30))
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось преобразовать Uid психолога в ObjectID | " + err.Error()}.SendMsg(w)
		return
	}
	newLogin := NewB64LowString(TrimStr(r.FormValue("login"), 40))
	newPas := Encrypt(TrimStr(r.FormValue("password"), 50))
	deleteFlag := r.FormValue("delete") == "true"
	newAvailable, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Некорректное доступное количество | " + err.Error()}.SendMsg(w)
	}
	var newTests []string
	for i := 1; i < TestsLen; i++ {
		if val, ok := TestDecode[r.FormValue("t"+strconv.Itoa(i))]; ok && val != "" {
			newTests = append(newTests, val)
		}
	}
	var curPsy Psy
	curPsyFilter := bson.M{"_id": psyUid, "status": PsyStatus}
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err = UsersCol.FindOne(ctx, curPsyFilter).Decode(&curPsy)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось получить текущие данные психолога | " + err.Error()}.SendMsg(w)
		return
	}

	newData := bson.M{"login": newLogin, "pas": newPas, "available": newAvailable, "tests": newTests}
	if curPsy.DeleteDate.IsZero() && deleteFlag {
		newData["deleteDate"] = time.Now()

		ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
		_, err = UsersCol.UpdateOne(ctx, curPsyFilter, bson.D{{"$set", newData}})
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: "Произошла ошибка при пре-удалении психолога | " + err.Error()}.SendMsg(w)
			return
		}

		go MoveTesteesByOwnerToArchive(psyUid.Hex())

	} else if !curPsy.DeleteDate.IsZero() && !deleteFlag {
		ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
		_, err = UsersCol.UpdateOne(ctx, curPsyFilter, bson.D{{"$set", newData}, {"$unset", bson.M{"deleteDate": 1}}})
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: "Произошла ошибка при восстановлении данных психолога | " + err.Error()}.SendMsg(w)
			return
		}

		go MoveTesteesByOwnerFromArchive(psyUid.Hex())

	} else {
		ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
		_, err = UsersCol.UpdateOne(ctx, curPsyFilter, bson.D{{"$set", newData}})
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: "Произошла ошибка при восстановлении данных психолога | " + err.Error()}.SendMsg(w)
			return
		}
	}

	JsonMsg{Kind: SucKind}.SendMsg(w)

})
