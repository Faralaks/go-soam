package tools

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func init() {
	var err error
	go ReadFeedBack()
	FeedBack = make(chan interface{})
	FeedBack <- "FeedBack is Ready!"

	folder := "."
	if runtime.GOOS != "darwin" {
		exec, _ := os.Executable()
		path := strings.Split(exec, "/")
		folder = strings.Join(path[:len(path)-1], "/")
	}

	configFile, err := ioutil.ReadFile(folder + "/config.json")
	if err != nil {
		VPrint("PANIC!!!  Ошибка при открытии конфига", err)
		panic(err)
	}

	err = json.Unmarshal(configFile, &configData)
	if err != nil {
		VPrint("PANIC!!!  Ошибка при декодировании JSON-конфига |" + err.Error())
		panic(err.Error())
	}

	atLifeTime, err := strconv.Atoi(configData["atLifeTime"])
	if err != nil {
		VPrint("Введено некорректное время жизни atLifeTime |", err.Error())
		panic(err)
	}
	Config.ATLifeTime = time.Duration(atLifeTime) * time.Minute

	rtLifeTime, err := (strconv.Atoi(configData["rtLifeTime"]))
	if err != nil {
		VPrint("Введено некорректное время жизни rtLifeTime |", err.Error())
		panic(err)
	}
	Config.RTLifeTime = time.Duration(rtLifeTime) * time.Minute

	Config.AccessSecret = []byte(configData["accessSecret"])
	Config.RefreshSecret = []byte(configData["refreshSecret"])
	Config.PasSecret = []byte(configData["pasSecret"])
	Config.Port = configData["port"]
	Config.Address = configData["address"]
	Config.CurPath = configData["curPath"]
	Config.MongoUrl = configData["mongoUrl"]
	Config.DbName = configData["dbName"]
	Config.UsersColName = configData["usersColName"]
	Config.TokensColName = configData["tokensColName"]
	FeedBack <- "Config is Ready!"

	client, err := mongo.NewClient(options.Client().ApplyURI(Config.MongoUrl))
	if err != nil {
		VPrint(err.Error())
		panic(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		VPrint(err.Error())
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		VPrint(err.Error())
		panic(err)
	}

	Client = client
	TokensCol = Client.Database(Config.DbName).Collection(Config.TokensColName)
	UsersCol = Client.Database(Config.DbName).Collection(Config.UsersColName)

	FeedBack <- "MongoDB is Ready!"

	c, err := aes.NewCipher(Config.PasSecret)
	if err != nil {
		VPrint(err.Error())
		panic(err)
	}
	Config.Gcm, err = cipher.NewGCM(c)
	if err != nil {
		VPrint(err.Error())
		panic(err)
	}
	FeedBack <- "Crypto is Ready!"

}
