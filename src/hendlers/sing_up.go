package hendlers

import (
	"context"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"time"
	. "tools"
)

var SignUp = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	ege, err := strconv.Atoi(r.FormValue("ege"))
	grade, err := strconv.Atoi(r.FormValue("grade"))

	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Недопустимое значение в одном из полей: Возраст, Класс | " + err.Error()}.SendMsg(w)
		return
	}

	newUser := FullUser{
		Uid:         p.NewObjectID(),
		Login:       NewB64LowString(TrimStr(r.FormValue("newLogin"), 40)),
		Pas:         NewSHA256(TrimStr(r.FormValue("newPassword"), 50)),
		Status:      TesteeStatus,
		CreatedDate: CurUtcStamp(),
		Step:        "",
		Name:        Encrypt(TrimStr(r.FormValue("name"), 20) + " " + TrimStr(r.FormValue("surname"), 30)),
		BirthYear:   uint16(time.Now().Year() - ege),
		Ege:         uint8(ege),
		Grade:       uint8(grade),
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	_, err = UsersCol.InsertOne(ctx, newUser)
	if err != nil {
		errCode := err.(mongo.WriteException).WriteErrors[0].Code
		if errCode == 11000 {
			JsonMsg{Kind: DuplicateKeyKind, Field: "newLogin"}.SendMsg(w)
		} else {
			JsonMsg{Kind: FatalKind, Msg: "Не удалось сохранить нового пользователя в базу данных | " + err.Error()}.SendMsg(w)
		}
		return
	}

	JsonMsg{Kind: SucKind}.SendMsg(w)

})
