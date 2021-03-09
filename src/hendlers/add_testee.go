package hendlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	. "tools"
)

type newGradesResponse struct {
	Kind      string  `json:"kind"`
	NewGrades *Grades `json:"newGrades"`
}

var Add_testees = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var curPsy Psy
	owner := r.Header.Get("owner")
	curPsyUid, err := p.ObjectIDFromHex(owner)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось преобразовать uid в ObjectID | " + err.Error()}.SendMsg(w)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = UsersCol.FindOne(ctx, bson.M{"_id": curPsyUid, "status": PsyStatus}).Decode(&curPsy)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось получить данные текущего психолога | " + err.Error()}.SendMsg(w)
		return
	}

	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Получено некорректное количество испытуемых | " + err.Error()}.SendMsg(w)
		return
	}
	if count > curPsy.Available {
		count = curPsy.Available
	}
	ident, err := curPsy.Ident.SensDecode()
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Почему-то не удалось превратить в строку идентификатор психолога | " + err.Error()}.SendMsg(w)
		return
	}

	grade := NewB64UpString(TrimStr(r.FormValue("grade"), 30))
	if grade == "" {
		JsonMsg{Kind: FatalKind, Msg: "Получено пустое название класса"}.SendMsg(w)
		return
	}

	var testees = make([]interface{}, count, count)
	for i := 0; i < count; i++ {
		rand.Seed(time.Now().UnixNano())
		testees[i] = bson.M{
			"_id":         p.NewObjectID(),
			"login":       NewB64LowString(ident + "_" + strconv.Itoa(curPsy.Counter+i)),
			"pas":         Encrypt(GeneratePas()),
			"owner":       owner,
			"status":      TesteeStatus,
			"createdDate": CurUtcStamp(),
			"tests":       curPsy.Tests,
			"msg":         "",
			"step":        "",
			"grade":       grade,
			"result":      rand.Intn(3) - 1,
		}
	}

	ctx, _ = context.WithTimeout(context.Background(), 4*time.Second)
	_, err = UsersCol.InsertMany(ctx, testees)
	if err != nil {
		if err.(mongo.WriteException).WriteErrors[0].Code == 11000 {
			JsonMsg{Kind: FatalKind, Msg: "Не удалось добавить  испытуемого в базу данных из за совпадения Логинов | " + err.Error()}.SendMsg(w)
		} else {
			JsonMsg{Kind: FatalKind, Msg: "Не удалось сохранить новых испытуемых в базу данных | " + err.Error()}.SendMsg(w)
		}
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	_, err = UsersCol.UpdateOne(ctx, bson.M{"_id": curPsyUid, "status": PsyStatus}, bson.D{
		{"$inc", bson.D{
			{fmt.Sprintf("grades.%v.whole", grade), count}, {fmt.Sprintf("grades.%v.not_yet", grade), count},
			{"counter", count}, {"available", -count},
		}},
	})
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось добавить  испытуемого в базу данных из за совпадения Логинов | " + err.Error()}.SendMsg(w)
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
	err = UsersCol.FindOne(ctx, bson.M{"_id": curPsyUid, "status": PsyStatus}).Decode(&curPsy)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось повторно получить данные текущего психолога | " + err.Error()}.SendMsg(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(newGradesResponse{SucKind, &curPsy.Grades})
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "В процессе конвертации в json возникла ошибка | " + err.Error()}.SendMsg(w)
		return
	}

})
