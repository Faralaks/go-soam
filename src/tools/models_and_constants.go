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
const TesteeStatus = "testee"

var AdminAccess = &[]string{AdminStatus}
var TesteeAccess = &[]string{TesteeStatus}

var AllAccess = &[]string{TesteeStatus, AdminStatus}

const NotYetResult = int8(-1)
const ClearResult = int8(0)
const DangerResult = int8(1)

const MskOffset = 10800

type User struct {
	Uid          p.ObjectID `json:"uid" bson:"_id"`
	Login        B64String  `json:"login" bson:"login"`
	VKToken      string     `json:"vk_token,omitempty" bson:"vk_token,omitempty"`
	Pas          string     `json:"pas,omitempty" bson:"pas,omitempty"`
	Status       string     `json:"status" bson:"status"`
	CreatedDate  Timestamp  `json:"create_date" bson:"createdDate"`
	ModifiedDate Timestamp  `json:"modifiedDate,omitempty" bson:"modifiedDate,omitempty"`
}

type FullUser struct {
	Uid          p.ObjectID             `json:"uid" bson:"_id"`
	Login        B64String              `json:"login" bson:"login"`
	Pas          string                 `json:"pas,omitempty" bson:"pas"`
	Status       string                 `json:"status" bson:"status"`
	CreatedDate  Timestamp              `json:"create_date" bson:"createdDate"`
	ModifiedDate Timestamp              `json:"modifiedDate,omitempty" bson:"modifiedDate,omitempty"`
	Step         uint8                  `json:"step,omitempty" bson:"step"`
	Tests        []string               `json:"tests,omitempty" bson:"tests"`
	Result       map[string]interface{} `json:"result,omitempty" bson:"result,omitempty"`
	Name         string                 `json:"name,omitempty" bson:"name"`
	BirthYear    uint16                 `json:"birthYear,omitempty" bson:"birthYear"`
	Ege          uint8                  `json:"ege,omitempty" bson:"ege"`
	Grade        uint8                  `json:"grade,omitempty" bson:"grade"`
}

type VKTokenData struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	UserId    int    `json:"user_id"`
	UserName  string `json:"userName,omitempty"`
}
