package hendlers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
	. "tools"
)

func RefreshMiddleware(cookieRt *http.Cookie, w http.ResponseWriter, r *http.Request, next http.Handler) {
	//println("~~~~~~~~~~~ RefreshMiddleware")
	rt := cookieRt.Value
	claims, err := ExtractRt(rt)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Недействительный ключ обновления | " + err.Error()}.SendMsg(w)
		return
	}
	var rtd RefreshTokenData
	uuid := claims["uuid"].(string)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//println(uuid)
	err = TokensCol.FindOne(ctx, bson.M{"_id": uuid}).Decode(&rtd)
	if err != nil {
		ctx, _ = context.WithTimeout(context.Background(), 7*time.Second)
		_, _ = TokensCol.DeleteMany(ctx, bson.M{"owner": claims["owner"].(string)})

		http.SetCookie(w, &http.Cookie{Name: "AccessToken", HttpOnly: true, MaxAge: -1})
		http.SetCookie(w, &http.Cookie{Name: "RefreshToken", HttpOnly: true, MaxAge: -1})

		JsonMsg{Kind: ReloginKind, Msg: "Скомпрометирован ключ обновления | " + err.Error()}.SendMsg(w)
		return
	}
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	_, err = TokensCol.DeleteOne(ctx, bson.M{"_id": uuid})
	if err != nil {
		println("\n", err.Error(), "\n\n")
	}
	NewAt, NewRt, err := CreateTokens(rtd.Owner, rtd.Status)
	if err != nil {
		JsonMsg{Kind: ReloginKind, Msg: "Не удалось создать новые токены | " + err.Error()}.SendMsg(w)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "AccessToken", Value: NewAt, HttpOnly: true, Expires: time.Now().Add(time.Minute * 10)})
	http.SetCookie(w, &http.Cookie{Name: "RefreshToken", Value: NewRt, HttpOnly: true, Expires: time.Now().Add(time.Hour)})

	r.Header.Set("owner", rtd.Owner)
	r.Header.Set("status", rtd.Status)

	next.ServeHTTP(w, r)
}
