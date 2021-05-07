package handlers

import (
	"net/http"
	"vk"
)

func Get_vk_token(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"][0]
	vk.GetUserOauthToken(code)

}
