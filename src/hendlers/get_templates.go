package hendlers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"net/http"
	"time"
	. "tools"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	main := Config.CurPath + "/templates/index.html"
	base := Config.CurPath + "/templates/base.html"

	tmpl, err := template.ParseFiles(main, base)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = tmpl.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	main := Config.CurPath + "/templates/admin.html"
	base := Config.CurPath + "/templates/base.html"

	tmpl, err := template.ParseFiles(main, base)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = tmpl.ExecuteTemplate(w, "admin", nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}

var TesteePage = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	curUserUid, err := p.ObjectIDFromHex(r.Header.Get("owner"))
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось преобразовать uid в ObjectID | " + err.Error()}.SendMsg(w)
		return
	}

	var user FullUser
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = UsersCol.FindOne(ctx, bson.M{"_id": curUserUid}).Decode(&user)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "Не удалось получить данные пользователя | " + err.Error()}.SendMsg(w)
		return
	}

	step := user.Step

	w.Header().Set("Content-Type", "text/html")
	main := Config.CurPath + "/templates/blank_" + step + ".html"
	base := Config.CurPath + "/templates/base.html"

	tmpl, err := template.ParseFiles(main, base)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = tmpl.ExecuteTemplate(w, "blank", nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

})
