package blank_handlers

import (
	"fmt"
	"strconv"
	"strings"
	p "go.mongodb.org/mongo-driver/bson/primitive"

)

var BlankHandlers = make(map[string]func([]string, p.ObjectID) error)
var Atoi = strconv.Atoi
var Itoa = strconv.Itoa

func init() {

		BlankHandlers["BPAQ"] = BPAQ_hendler
		BlankHandlers["ITO"] = ITO_hendler
}

func intValMap(ansList []string) (map[string]int, error) {
	ans := make(map[string]int)
	for i := 0; i < len(ansList); i++ {
		tmp := strings.Split(ansList[i], "=")
		val, err := Atoi(tmp[1])
		if err != nil {
			return nil, fmt.Errorf("%v элемент не удалось приорзовать в число | %v", tmp[0], err.Error())
		}
		ans[tmp[0]] = val
	}
	return ans, nil
}


func qNum(num int) string {
	return "q" + Itoa(num)
}
