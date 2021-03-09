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

var Del_result = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var curTestee Testee
	testeeUidStr := TrimStr(r.FormValue("testeeUid"), 30)
	//fmt.Printf("\n\ntesteeUidStr:%v\n", testeeUidStr)
	testeeUid, err := p.ObjectIDFromHex(testeeUidStr)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось преобразовать uid испытуемого в ObjectID | " + err.Error()}.SendMsg(w)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = UsersCol.FindOne(ctx, bson.M{"_id": testeeUid, "status": TesteeStatus}).Decode(&curTestee)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось получить данные испытуемого | " + err.Error()}.SendMsg(w)
		return
	}

	grade := B64Enc(TrimStr(r.FormValue("grade"), 30))
	if grade == "" {
		JsonMsg{Kind: FatalKind, Msg: "Получено некорректное значение для класса испытуемых"}.SendMsg(w)
		return
	}
	//fmt.Printf("msg:%v\nTimMsg:%v\n\n", r.FormValue("msg"), TrimStr(r.FormValue("msg"), 500))
	msg := B64Enc(TrimStr(r.FormValue("msg"), 500))
	if msg == "" {
		JsonMsg{Kind: FatalKind, Msg: "Получено некорректное значение для причины удаления "}.SendMsg(w)
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 4*time.Second)
	_, err = UsersCol.UpdateOne(ctx, bson.M{"_id": testeeUid, "status": TesteeStatus}, bson.D{{"$set", bson.D{{"msg", msg}, {"result", NotYetResult}}}})
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось изменить результат и установить причину удаления у испытуемого| " + err.Error()}.SendMsg(w)
		return
	}

	curPsyUid, err := p.ObjectIDFromHex(curTestee.Owner)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Почему-то не удалось преобразовать Uid хозяина в  ObjectID |" + err.Error()}.SendMsg(w)
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 4*time.Second)
	_, err = UsersCol.UpdateOne(ctx, bson.M{"_id": curPsyUid, "status": PsyStatus}, bson.D{
		{"$inc", bson.D{
			{fmt.Sprintf("grades.%v.msg", grade), 1}, {fmt.Sprintf("grades.%v.not_yet", grade), 1},
			{fmt.Sprintf("grades.%v.%v", grade, ResultDecode[curTestee.Result]), -1},
		}},
	})
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Произошла ошибка при изменении данных класса у психолога | " + err.Error()}.SendMsg(w)
		return
	}

})
