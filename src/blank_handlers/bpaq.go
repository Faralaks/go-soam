package blank_handlers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	. "tools"
)

func BPAQ_hendler(ansList []string, uid p.ObjectID) error {
	ans, err := intValMap(ansList)
	if err != nil {
		return err
	}

	result := map[string]int{"aggression": 0, "anger": 0, "hostility": 0}

	aggressionList := [9]string{"1", "4", "7", "10", "13", "16", "22", "24"}
	angerList := [7]string{"2", "5", "8", "14", "17", "20"}
	hostilityList := [8]string{"3", "6", "9", "12", "15", "18", "21", "23"}

	// Считаем по шкале Физическая агрессия
	for _, n := range aggressionList {
		result["aggression"] += ans["q"+n]
	}
	result["aggression"] += (6 - ans["q19"])

	// Счетаем шкалу Гнева
	for _, n := range angerList {
		result["anger"] += ans["q"+n]
	}
	result["anger"] += (6 - ans["q11"])

	// Счетаем шкалу Враждебности
	for _, n := range hostilityList {
		result["hostility"] += ans["q"+n]
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	_, err = UsersCol.UpdateOne(ctx, bson.M{"_id": uid}, bson.D{{"$set", bson.M{"result.BPAQ": result}}, {"$inc", bson.M{"step": 1}}})
	if err != nil {
		return err
	}

	ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
	_, err = ResultsCol.InsertOne(ctx, bson.M{"pastedBy": uid, "date": CurUtcStamp(), "answers": ans})
	if err != nil {
		return err
	}

	return nil
}
