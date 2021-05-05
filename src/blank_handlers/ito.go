package blank_handlers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	. "tools"
)

func ITO_hendler(ansList []string, uid p.ObjectID) error {
	ans, err := intValMap(ansList)
	if err != nil {
		return err
	}

	result := map[string]int{"L": 0,"F": 0,"i": 0,"ii": 0,"iii": 0,"iv": 0,"v": 0,"vi": 0,"vii": 0,"viii": 0}

	LTList := [9]string{"16", "31", "45", "46", "60", "61", "75", "76", "90"}

	FTList := [9]string{"2", "17", "32", "47", "62", "77", "64", "79"}

	iTList := [6]string{"12", "27", "29", "42", "44", "72"}
	iFList := [3]string{"14", "57", "87"}

	iiTList := [6]string{"4", "19", "21", "34", "49", "50"}
	iiFList := [3]string{"6", "65", "80"}

	iiiTList := [7]string{"7", "22", "36", "37", "51", "53", "68"}
	iiiFList := [2]string{"66", "81"}

	ivTList := [6]string{"9", "24", "26", "39", "41", "56"}
	ivFList := [3]string{"71", "83", "86"}

	vTList := [6]string{"3", "5", "33", "35", "48", "78"}
	vFList := [3]string{"18", "20", "63"}

	viTList := [5]string{"15", "28", "43", "59", "89"}
	viFList := [4]string{"11", "13", "30", "74"}

	viiTList := [7]string{"8", "23", "38", "52", "54", "69", "84"}
	viiFList := [2]string{"67", "82"}

	viiiTList := [5]string{"10", "25", "40", "55", "58"}
	viiiFList := [4]string{"70", "73", "85", "88"}



	// Шкала Ложь
	for _, n := range LTList {
		if ans["q"+n] == 1 { result["L"] += 1 }
	}
	// Шкала Аггравация
	for _, n := range FTList {
		if ans["q"+n] == 1 { result["F"] += 1 }
	}


	// Шкала Экстраверсия Верно
	for _, n := range iTList {
		if ans["q"+n] == 1 { result["i"] += 1 }
	}
	// Шкала Экстраверсия Неерно
	for _, n := range iFList {
		if ans["q"+n] == 2 { result["i"] += 1 }
	}


	// Шкала Спонтанность Верно
	for _, n := range iiTList {
		if ans["q"+n] == 1 { result["ii"] += 1 }
	}
	// Шкала Спонтанность Неерно
	for _, n := range iiFList {
		if ans["q"+n] == 2 { result["ii"] += 1 }
	}


	// Шкала Агрессивность Верно
	for _, n := range iiiTList {
		if ans["q"+n] == 1 { result["iii"] += 1 }
	}
	// Шкала Агрессивность Неерно
	for _, n := range iiiFList {
		if ans["q"+n] == 2 { result["iii"] += 1 }
	}


	// Шкала Ригидность Верно
	for _, n := range ivTList {
		if ans["q"+n] == 1 { result["iv"] += 1 }
	}
	// Шкала Ригидность Неерно
	for _, n := range ivFList {
		if ans["q"+n] == 2 { result["iv"] += 1 }
	}


	// Шкала Интроверсия Верно
	for _, n := range vTList {
		if ans["q"+n] == 1 { result["v"] += 1 }
	}
	// Шкала Интроверсия Неерно
	for _, n := range vFList {
		if ans["q"+n] == 2 { result["v"] += 1 }
	}


	// Шкала Сензитивность Верно
	for _, n := range viTList {
		if ans["q"+n] == 1 { result["vi"] += 1 }
	}
	// Шкала Сензитивность Неерно
	for _, n := range viFList {
		if ans["q"+n] == 2 { result["vi"] += 1 }
	}


	// Шкала Тревожность Верно
	for _, n := range viiTList {
		if ans["q"+n] == 1 { result["vii"] += 1 }
	}
	// Шкала Тревожность Неерно
	for _, n := range viiFList {
		if ans["q"+n] == 2 { result["vii"] += 1 }
	}


	// Шкала Эмотивность Верно
	for _, n := range viiiTList {
		if ans["q"+n] == 1 { result["viii"] += 1 }
	}
	// Шкала Эмотивность Неерно
	for _, n := range viiiFList {
		if ans["q"+n] == 2 { result["viii"] += 1 }
	}


	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	_, err = UsersCol.UpdateOne(ctx, bson.M{"_id": uid}, bson.D{{"$set", bson.M{"result.ITO": result}}, {"$inc", bson.M{"step": 1}}})
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
