package tools

import (
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type B64String string

func NewB64String(l string) B64String {
	return B64String(B64Enc(l))
}
func NewB64LowString(l string) B64String {
	return B64String(B64Enc(strings.ToLower(l)))
}
func NewB64UpString(l string) B64String {
	return B64String(B64Enc(strings.ToUpper(l)))
}

func (s B64String) SensDecode() (string, error) {
	return B64Dec(string(s))
}

func (s B64String) Decode() string {
	text, err := B64Dec(string(s))
	if err != nil {
		return "Err: " + err.Error()
	}
	return text
}

type Timestamp int64

func NewTimestamp(d time.Time) Timestamp {
	return Timestamp(d.UTC().Unix())
}
func CurUtcStamp() Timestamp {
	return NewTimestamp(time.Now())
}

const AdminStatus = "admin"
const PsyStatus = "psy"
const TesteeStatus = "testee"

var AdminAccess = &[]string{AdminStatus}
var PsyAccess = &[]string{PsyStatus}

//var TesteeAccess = []string{TesteeStatus}
var AdminAndPsyAccess = &[]string{PsyStatus, AdminStatus}
var AllAccess = &[]string{PsyStatus, TesteeStatus, AdminStatus}

const NotYetResult = int8(-1)
const ClearResult = int8(0)
const DangerResult = int8(1)

var ResultDecode = map[int8][2]string{
	NotYetResult: {"not_yet", "Нет результата"},
	ClearResult:  {"clear", "Вне группы"},
	DangerResult: {"danger", "В группе"},
}

var TestDecode = map[string]string{
	"1": "1",
	"2": "2",
}

var TestsLen = len(ResultDecode)

const MskOffset = 10800

type Grades map[B64String]struct {
	Whole   uint `json:"whole,omitempty" bson:"whole,omitempty"`
	Not_yet uint `json:"not_yet,omitempty" bson:"not_yet,omitempty"`
	Clear   uint `json:"clear,omitempty" bson:"clear,omitempty"`
	Danger  uint `json:"danger,omitempty" bson:"danger,omitempty"`
	Msg     uint `json:"msg,omitempty" bson:"msg,omitempty"`
}

type User struct {
	Uid          p.ObjectID `json:"uid" bson:"_id"`
	Login        B64String  `json:"login" bson:"login"`
	Pas          string     `json:"pas" bson:"pas"`
	Status       string     `json:"status" bson:"status"`
	CreatedDate  Timestamp  `json:"create_date" bson:"createdDate"`
	Owner        string     `json:"added_by" bson:"owner"`
	ModifiedDate Timestamp  `json:"modifiedDate" bson:"modifiedDate"`
}

type Psy struct {
	Uid          p.ObjectID `json:"uid" bson:"_id"`
	Login        B64String  `json:"login" bson:"login"`
	Pas          string     `json:"pas" bson:"pas"`
	Status       string     `json:"status" bson:"status"`
	CreatedDate  Timestamp  `json:"create_date" bson:"createdDate"`
	Ident        B64String  `json:"ident" bson:"ident"`
	Owner        string     `json:"added_by" bson:"owner"`
	Available    int        `json:"count" bson:"available"`
	Tests        []string   `json:"tests" bson:"tests"`
	Grades       Grades     `json:"grades" bson:"grades"`
	Counter      int        `bson:"counter,omitempty"`
	DeleteDate   time.Time  `json:"pre_del,omitempty" bson:"deleteDate,omitempty"`
	ModifiedDate Timestamp  `json:"modifiedDate" bson:"modifiedDate"`
}
type Testee struct {
	Uid         p.ObjectID `json:"uid" bson:"_id"`
	Login       B64String  `json:"login" bson:"login"`
	Pas         string     `json:"pas" bson:"pas"`
	Owner       string     `json:"added_by" bson:"owner"`
	Status      string     `json:"status" bson:"status"`
	CreatedDate Timestamp  `json:"create_date" bson:"createdDate"`
	Tests       []string   `json:"tests" bson:"tests"`
	Msg         B64String  `json:"msg" bson:"msg"`
	Grade       B64String  `json:"grade" bson:"grade"`
	Step        string     `json:"step" bson:"step"`
	Result      int8       `json:"result" bson:"result"`
}

type MultiUser struct {
	Uid          p.ObjectID `json:"uid" bson:"_id"`
	Login        B64String  `json:"login" bson:"login"`
	Pas          string     `json:"pas" bson:"pas"`
	Status       string     `json:"status" bson:"status"`
	CreatedDate  Timestamp  `json:"create_date" bson:"createdDate"`
	DeleteDate   time.Time  `json:"pre_del,omitempty" bson:"deleteDate,omitempty"`
	Ident        B64String  `json:"ident,omitempty" bson:"ident,omitempty"`
	Available    int        `json:"count" bson:"available"`
	Owner        string     `json:"added_by" bson:"owner"`
	Tests        []string   `json:"tests,omitempty" bson:"tests,omitempty"`
	Grades       Grades     `json:"grades,omitempty" bson:"grades,omitempty"`
	Msg          B64String  `json:"msg,omitempty" bson:"msg,omitempty"`
	Grade        B64String  `json:"grade" bson:"grade"`
	Step         string     `json:"step,omitempty" bson:"step,omitempty"`
	Result       int8       `json:"result,omitempty" bson:"result,omitempty"`
	Counter      int        `json:"counter,omitempty" bson:"counter,omitempty"`
	ModifiedDate Timestamp  `json:"modifiedDate" bson:"modifiedDate"`
}
