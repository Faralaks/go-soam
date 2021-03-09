package hendlers

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	. "tools"
)

var Accept_del = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var delTestee Testee
	testeeUidStr := TrimStr(r.FormValue("testeeUid"), 30)
	//fmt.Printf("\n\ntesteeUidStr:%v\n", testeeUidStr)
	testeeUid, err := p.ObjectIDFromHex(testeeUidStr)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось преобразовать uid испытуемого в ObjectID | " + err.Error()}.SendMsg(w)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err = UsersCol.FindOne(ctx, bson.M{"_id": testeeUid, "status": TesteeStatus}).Decode(&delTestee)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось получить данные испытуемого | " + err.Error()}.SendMsg(w)
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
	_, err = UsersCol.UpdateOne(ctx, bson.M{"_id": testeeUid, "status": TesteeStatus}, bson.D{{"$set", bson.D{{"msg", ""}}}})
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось изменить причину удаления у испытуемого| " + err.Error()}.SendMsg(w)
		return
	}

	curPsyUid, err := p.ObjectIDFromHex(delTestee.Owner)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Почему-то не удалось преобразовать Uid хозяина в  ObjectID |" + err.Error()}.SendMsg(w)
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
	_, err = UsersCol.UpdateOne(ctx, bson.M{"_id": curPsyUid, "status": PsyStatus}, bson.D{
		{"$inc", bson.D{
			{fmt.Sprintf("grades.%v.msg", delTestee.Grade), -1}, {"available", 1},
		}},
	})
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Произошла ошибка при изменении данных психолога после подтверждения ужаления | " + err.Error()}.SendMsg(w)
		return
	}

})
