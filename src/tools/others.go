package tools

import (
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type configType struct {
	CurPath           string
	Port              string
	Address           string
	Gcm               cipher.AEAD
	PasSecret         []byte
	AccessSecret      []byte
	RefreshSecret     []byte
	ATLifeTime        time.Duration
	RTLifeTime        time.Duration
	MongoUrl          string
	DbName            string
	UsersColName      string
	TokensColName     string
	ResultsColName    string
	StatsColName      string
	TestList          []string
	TestCount         uint8
	OauthClientID     string
	OauthKey          string
	OauthRedirectPath string
	OauthRedirectURL  string
}

var configData map[string]string
var Config = configType{}

var Client *mongo.Client
var TokensCol *mongo.Collection
var UsersCol *mongo.Collection
var ResultsCol *mongo.Collection

var FeedBack chan interface{}
var feedCounter = 0

func ReadFeedBack() {
	for {
		fmt.Printf("%v\n", <-FeedBack)
	}
}

func VPrint(lines ...interface{}) {
	fmt.Printf("\n======= %v   ¯\\_(ツ)_/¯    =======\n", feedCounter)
	for i, s := range lines {
		if fmt.Sprintf("%T", s) == "string" {
			fmt.Printf("%v %T\t\"%v\" l: %v\n", i, s, s, len(s.(string)))
		} else {
			fmt.Printf("%v %T\t%v\n", i, s, s)
		}
	}
	fmt.Print("=======   ~~~~~~~~~~~~    =======\n")
	feedCounter++

}

const FatalKind = "Fatal"
const ReloginKind = "Relogin"
const BadAuthKind = "BadAuth"
const SucKind = "Suc"
const BadUpdateKind = "BadUpdate"
const DuplicateKeyKind = "DuplicatedField"

type JsonMsg struct {
	Kind  string `json:"kind,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Field string `json:"field,omitempty"`
}

func (jm JsonMsg) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(jm)
	if jm.Kind != SucKind {
		log.Printf("ErrKind: %v, Msg: %v, Field: %v", jm.Kind, jm.Msg, jm.Field)
	}
}

func _(elem string, list []string) bool {
	for _, e := range list {
		if elem == e {
			return true
		}
	}
	return false
}

func TrimStr(str string, l int) string {
	str = strings.TrimSpace(str)
	runes := []rune(str)
	if len(runes) > l {
		str = string(runes[:l])
	}
	return str
}

func GeneratePas() string {
	return "pas"
}

func IsAllowed(status string, allowList *[]string) bool {
	for _, allow := range *allowList {
		if status == allow {
			return true
		}
	}
	return false

}
func SetUserDataHeaders(status, owner string, r *http.Request) {
	r.Header.Set("status", status)
	r.Header.Set("owner", owner)
}

func SetLoginCookies(w http.ResponseWriter, newAt, newRt string) {
	http.SetCookie(w, &http.Cookie{Name: "AccessToken", Value: newAt, HttpOnly: true, Expires: time.Now().UTC().Add(Config.ATLifeTime)})
	http.SetCookie(w, &http.Cookie{Name: "RefreshToken", Value: newRt, HttpOnly: true, Expires: time.Now().UTC().Add(Config.RTLifeTime)})
}

func DeleteLoginCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: "AccessToken", HttpOnly: true, MaxAge: -1})
	http.SetCookie(w, &http.Cookie{Name: "RefreshToken", HttpOnly: true, MaxAge: -1})
}

func OpenInBrowser(url string) {
	if runtime.GOOS == "darwin" { // macOS
		_ = exec.Command("open", url).Start()
	}

}
