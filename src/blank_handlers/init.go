package blank_handlers

import "strconv"

var BlankHandlers = make(map[string]func(map[string]string) error)

func init() {
	BlankHandlers["BPAQ"] = BPAQ_hendler
}

func qNum(num int) string {
	return "q" + strconv.Itoa(num)
}
