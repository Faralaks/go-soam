package hendlers

import (
	"context"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"strings"
	"time"
	. "tools"
)

var Add_psy = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	available, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Некорректное доступное количество | " + err.Error()}.SendMsg(w)
	}
	var tests []string
	for i := 1; i <= TestsLen; i++ {
		if val, ok := TestDecode[r.FormValue("t"+strconv.Itoa(i))]; ok && val != "" {
			tests = append(tests, val)
		}
	}
	newPsy := Psy{
		Uid:          p.NewObjectID(),
		Login:        NewB64LowString(TrimStr(r.FormValue("login"), 40)),
		Pas:          Encrypt(TrimStr(r.FormValue("password"), 50)),
		Status:       PsyStatus,
		CreatedDate:  CurUtcStamp(),
		Ident:        NewB64LowString(TrimStr(r.FormValue("ident"), 40)),
		Owner:        r.Header.Get("owner"),
		Available:    available,
		Tests:        tests,
		Grades:       Grades{},
		ModifiedDate: CurUtcStamp(),
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	_, err = UsersCol.InsertOne(ctx, newPsy)
	if err != nil {
		errCode := err.(mongo.WriteException).WriteErrors[0].Code
		if errCode == 11000 {
			field := strings.Split(err.Error(), "{")[3][1:6]
			JsonMsg{Kind: DuplicateKeyKind, Field: field}.SendMsg(w)
		} else {
			JsonMsg{Kind: FatalKind, Msg: "Не удалось сохранить нового психолога в базу данных | " + err.Error()}.SendMsg(w)
		}
		return
	}

	JsonMsg{Kind: SucKind}.SendMsg(w)

})
