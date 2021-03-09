package hendlers

import (
	"html/template"
	"net/http"
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
	stats := Config.CurPath + "/templates/stats.html"
	base := Config.CurPath + "/templates/base.html"

	tmpl, err := template.ParseFiles(main, base, stats)
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

func PsyPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	main := Config.CurPath + "/templates/psy.html"
	stats := Config.CurPath + "/templates/stats.html"
	base := Config.CurPath + "/templates/base.html"

	tmpl, err := template.ParseFiles(main, base, stats)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = tmpl.ExecuteTemplate(w, "psy", nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}
