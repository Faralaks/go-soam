package hendlers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	. "tools"
)

var Edit_user_data = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	owner := r.Header.Get("owner")
	userUid, err := p.ObjectIDFromHex(owner)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось преобразовать Uid пользователя в ObjectID | " + err.Error()}.SendMsg(w)
		return
	}
	userStatus := r.Header.Get("status")

	oldPas := Encrypt(TrimStr(r.FormValue("password"), 50))
	newPas := Encrypt(TrimStr(r.FormValue("newPassword"), 50))

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	result, err := UsersCol.UpdateOne(ctx, bson.M{"_id": userUid, "status": userStatus, "pas": oldPas}, bson.D{{"$set", bson.M{"modifiedDate": CurUtcStamp(), "pas": newPas}}})
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Произошла ошибка при смене данных пользователя | " + err.Error()}.SendMsg(w)
		return
	}

	if result.MatchedCount == 1 {
		ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
		_, err = TokensCol.DeleteMany(ctx, bson.M{"owner": owner})
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: "Произошла ошибка при удалении токенов пользователя | " + err.Error()}.SendMsg(w)
			return
		}

		JsonMsg{Kind: ReloginKind}.SendMsg(w)
		return
	}

	JsonMsg{Kind: BadUpdateKind}.SendMsg(w)

})
