package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	. "hendlers"
	"net/http"
	"os"
	. "tools"
)

func main() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir(Config.CurPath + "/public"))
	r.PathPrefix("/js/").Handler(fs)
	r.Path("/favicon.ico").Handler(fs)

	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/login", Login).Methods("POST")

	r.HandleFunc("/admin", AdminPage).Methods("GET")
	r.HandleFunc("/psy", PsyPage).Methods("GET")

	r.Handle("/get_psy_list", AuthMiddleware(Get_psy_list, AdminAccess)).Methods("GET")
	r.Handle("/add_psy", AuthMiddleware(Add_psy, AdminAccess)).Methods("POST")
	r.Handle("/accept_del", AuthMiddleware(Accept_del, AdminAccess)).Methods("POST")
	r.Handle("/edit_psy", AuthMiddleware(Edit_psy, AdminAndPsyAccess)).Methods("POST")
	r.Handle("/get_user_data", AuthMiddleware(Get_user_data, AllAccess)).Methods("GET")
	r.Handle("/add_testees", AuthMiddleware(Add_testees, PsyAccess)).Methods("POST")
	r.Handle("/get_testee_list", AuthMiddleware(Get_testee_list, AdminAndPsyAccess)).Methods("GET")
	r.Handle("/del_result", AuthMiddleware(Del_result, PsyAccess)).Methods("POST")
	r.Handle("/edit_user_data", AuthMiddleware(Edit_user_data, AdminAndPsyAccess)).Methods("POST")
	r.Handle("/download", AuthMiddleware(Download, AdminAndPsyAccess)).Methods("GET")

	r.HandleFunc("/remake", remakeDb).Methods("GET")
	r.HandleFunc("/logout", logOut).Methods("GET")

	OpenInBrowser("http://127.0.0.1:" + Config.Port)
	_ = http.ListenAndServe(Config.Address+":"+Config.Port, handlers.LoggingHandler(os.Stdout, r))
}

var logOut = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	DeleteLoginCookies(w)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.Redirect(w, r, "/", 301)

})

var remakeDb = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	UsersCol.DeleteMany(context.TODO(), bson.M{})
	TokensCol.DeleteMany(context.TODO(), bson.M{})
	u := User{
		Uid:          p.NewObjectID(),
		Login:        NewB64String("master"),
		Pas:          Encrypt("retsam"),
		Status:       AdminStatus,
		CreatedDate:  CurUtcStamp(),
		Owner:        "Faralaks",
		ModifiedDate: CurUtcStamp(),
	}
	UsersCol.InsertOne(context.TODO(), u)

	_ = Psy{
		Uid:          p.NewObjectID(),
		Login:        NewB64LowString("psy"),
		Pas:          Encrypt(""),
		Status:       PsyStatus,
		CreatedDate:  CurUtcStamp(),
		Ident:        NewB64LowString("id"),
		Owner:        "Faralaks",
		Available:    9999,
		Tests:        []string{"1", "2"},
		Grades:       Grades{},
		ModifiedDate: CurUtcStamp(),
	}

	//UsersCol.InsertOne(context.TODO(), p)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.Redirect(w, r, "/", http.StatusMovedPermanently)

})
