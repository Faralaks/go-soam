package hendlers

import (
	"net/http"
	. "tools"
)

func AuthMiddleware(next http.Handler, allowList *[]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieAt, err := r.Cookie("AccessToken")
		if err != nil || len(cookieAt.Value) == 0 {
			if cookieRt, err := r.Cookie("RefreshToken"); err == nil && len(cookieRt.Value) != 0 {
				RefreshMiddleware(cookieRt, allowList, w, r, next)
				return
			}
			JsonMsg{Kind: ReloginKind, Msg: "Не были получены необходимые ключи"}.SendMsg(w)
			return
		}
		at := cookieAt.Value
		claims, err := ExtractToken(at, Config.AccessSecret)
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: "Недействительный ключ авторизации | " + err.Error()}.SendMsg(w)
			return
		}

		if !IsAllowed(claims["status"].(string), allowList) {
			JsonMsg{Kind: ReloginKind, Msg: "Отказано в доступе"}.SendMsg(w)
			return
		}

		SetUserDataHeaders(claims["status"].(string), claims["owner"].(string), r)

		next.ServeHTTP(w, r)

	})
}
