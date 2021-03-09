package hendlers

import (
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
	"strings"
	"time"
	. "tools"
)

var Download = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	filter := bson.M{"status": TesteeStatus}
	filename := "file.xlsx"
	owner := r.Header.Get("owner")
	if r.Header.Get("status") == AdminStatus {
		owner = r.FormValue("psyUid")
		if owner != "" {
			filter["owner"] = TrimStr(owner, 30)
		}
	}

	grade := TrimStr(r.FormValue("grade"), 30)
	if grade != "" {
		filter["grade"] = NewB64UpString(grade)
	}
	var file *excelize.File
	switch TrimStr(r.FormValue("target"), 20) {
	case "not_yet":
		filter["result"] = -1
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		cur, err := UsersCol.Find(ctx, filter)
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: "Не удалось извлечь список испытуемых | " + err.Error()}.SendMsg(w)
			return
		}
		if cur == nil {
			JsonMsg{Kind: FatalKind, Msg: "ТАКОГО НЕ МОЖЕТ СЛУЧИТСЯ, ПОТОМУ ЧТО НЕ МОЖЕТ. КУРСОР, МАТЬ ЕГО, СТАЛ NIL"}.SendMsg(w)
			return
		}
		defer cur.Close(context.TODO())

		file, err = excelize.OpenFile(Config.CurPath + "/not_yet_template.xlsx")
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: "Не удалось открыть заготовку not_yet_template.xlsx | " + err.Error()}.SendMsg(w)
			return
		}
		row := 2
		ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
		for cur.Next(ctx) {
			var elem Testee
			err := cur.Decode(&elem)
			if err != nil {
				JsonMsg{Kind: FatalKind, Msg: "В процессе декодирования испытуемых произошла ошибка | " + err.Error()}.SendMsg(w)
				return
			}
			pas, err := Decrypt(elem.Pas)
			login := elem.Login.Decode()
			if err != nil {
				JsonMsg{Kind: FatalKind, Msg: fmt.Sprintf("В процессе дешифрования пароля пользователя %v произошла ошибка | %v", elem.Login.Decode(), err.Error())}.SendMsg(w)
				return
			}
			_ = file.SetCellValue("Логины", "A"+strconv.Itoa(row), login)
			_ = file.SetCellValue("Логины", "B"+strconv.Itoa(row), pas)
			VPrint(elem.Msg)
			if elem.Msg != "" {
				_ = file.SetCellValue("Логины", "D"+strconv.Itoa(row), "Удален")
			}
			row++
		}
		if err := cur.Err(); err != nil {
			JsonMsg{Kind: FatalKind, Msg: "У курсора произошла ошибка | " + err.Error()}.SendMsg(w)
			return
		}
		filename = "Логины.xlsx"

	case "done":
		filter["result"] = bson.D{{"$gt", -1}}
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		cur, err := UsersCol.Find(ctx, filter)
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: "Не удалось извлечь список испытуемых | " + err.Error()}.SendMsg(w)
			return
		}
		if cur == nil {
			JsonMsg{Kind: FatalKind, Msg: "ТАКОГО НЕ МОЖЕТ СЛУЧИТСЯ, ПОТОМУ ЧТО НЕ МОЖЕТ. КУРСОР, МАТЬ ЕГО, СТАЛ NIL"}.SendMsg(w)
			return
		}
		defer cur.Close(context.TODO())

		file, err = excelize.OpenFile(Config.CurPath + "/done_template.xlsx")
		if err != nil {
			JsonMsg{Kind: FatalKind, Msg: "Не удалось открыть заготовку done_template.xlsx | " + err.Error()}.SendMsg(w)
			return
		}
		row := 2
		ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
		for cur.Next(ctx) {
			var elem Testee
			var rowStr = strconv.Itoa(row)
			err := cur.Decode(&elem)
			if err != nil {
				JsonMsg{Kind: FatalKind, Msg: "В процессе декодирования испытуемых произошла ошибка | " + err.Error()}.SendMsg(w)
				return
			}

			_ = file.SetCellValue("Результаты", "A"+rowStr, elem.Login.Decode())
			_ = file.SetCellValue("Результаты", "B"+rowStr, elem.Grade.Decode())
			_ = file.SetCellValue("Результаты", "C"+rowStr, strings.Join(elem.Tests, ", "))
			_ = file.SetCellValue("Результаты", "D"+rowStr, time.Unix(int64(elem.CreatedDate), 0).UTC().Add(MskOffset*time.Second).Format("2006-01-02 15:04:05")+" МСК")
			_ = file.SetCellValue("Результаты", "E"+rowStr, elem.Result)
			_ = file.SetCellValue("Результаты", "F"+rowStr, ResultDecode[elem.Result][1])
			row++
		}
		if err := cur.Err(); err != nil {
			JsonMsg{Kind: FatalKind, Msg: "У курсора произошла ошибка | " + err.Error()}.SendMsg(w)
			return
		}
		for row := 2; row < 9; row++ {
			strRow := strconv.Itoa(row)
			formula, _ := file.GetCellFormula("Результаты", "J"+strRow)
			_ = file.SetCellFormula("Результаты", "I"+strRow, formula)
			_ = file.SetCellFormula("Результаты", "J"+strRow, "")
			_ = file.SetCellStr("Результаты", "J"+strRow, "")
		}
		filename = "Результаты.xlsx"
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, filename))
	w.Header().Set("File-Name", filename)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	err := file.Write(w)
	if err != nil {
		JsonMsg{Kind: FatalKind, Msg: "В процессе конвертации в json возникла ошибка | " + err.Error()}.SendMsg(w)
		return
	}
})
