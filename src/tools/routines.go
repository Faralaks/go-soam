package tools

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func MoveTesteesByOwnerToArchive(owner string) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := UsersCol.UpdateMany(ctx, bson.M{"owner": owner, "status": TesteeStatus}, bson.D{
		{"$set", bson.D{{"deleteDate", time.Now()}}},
	})
	if err != nil {
		go VPrint(fmt.Sprintf("Произошла ошибка при установки метки deleteDate у испытуемых с хозяином %v | %v", owner, err.Error()))
	}
	time.Sleep(10 * time.Second)
	go VPrint(fmt.Sprintf("Испытуемые хозяина %v перемещены в архив", owner))
}

func MoveTesteesByOwnerFromArchive(owner string) {
	time.Sleep(20 * time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := UsersCol.UpdateMany(ctx, bson.M{"owner": owner, "status": TesteeStatus}, bson.D{
		{"$unset", bson.D{{"deleteDate", 1}}},
	})
	if err != nil {
		go VPrint(fmt.Sprintf("Произошла ошибка при снятии метки deleteDate у испытуемых с хозяином %v | %v", owner, err.Error()))

	}
	go VPrint(fmt.Sprintf("Испытуемые хозяина %v возвращены из архива", owner))
}
