package hendlers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
	. "tools"
)

func RefreshMiddleware(cookieRt *http.Cookie, allowList *[]string, w http.ResponseWriter, r *http.Request, next http.Handler) {
	rt := cookieRt.Value
	claims, err := ExtractToken(rt, Config.RefreshSecret)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Недействительный ключ обновления | " + err.Error()}.SendMsg(w)
		return
	}
	var rtd TokenData
	uid := claims["uid"].(string)
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = TokensCol.FindOne(ctx, bson.M{"_id": uid}).Decode(&rtd)
	if err != nil {
		ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
		_, _ = TokensCol.DeleteMany(ctx, bson.M{"owner": claims["owner"].(string)})

		DeleteLoginCookies(w)

		JsonMsg{Kind: ReloginKind, Msg: "Веронятно, ваш ключ обновления скомпромитирован, рекомендуется сменить пароль | " + err.Error()}.SendMsg(w)
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
	_, err = TokensCol.DeleteOne(ctx, bson.M{"_id": uid})
	if err != nil {
		VPrint(err.Error())
	}
	newAt, newRt, err := CreateTokens(rtd.Owner, rtd.Status)
	if err != nil {
		JsonMsg{Kind: ReloginKind, Msg: "Не удалось создать новые токены | " + err.Error()}.SendMsg(w)
		return
	}
	SetLoginCookies(w, newAt, newRt)

	if !IsAllowed(rtd.Status, allowList) {
		JsonMsg{Kind: ReloginKind, Msg: "Отказано в доступе"}.SendMsg(w)
		return
	}

	SetUserDataHeaders(claims["status"].(string), claims["owner"].(string), r)

	next.ServeHTTP(w, r)

}
