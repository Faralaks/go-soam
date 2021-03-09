let psyList, gradeList, testeeList;
let curPsy, curGrade;
let gradeCounters = {};
let fullCounter;
let needToReload = false;
let myData;

function setToDefault() {
    curPsy.dec_login = title(b64dec(curPsy.login));
    curPsy.dec_ident = title(b64dec(curPsy.ident));
    if (!curPsy.tests) curPsy.tests = [];
    $("#psyFormBtnSave").prop("disabled", true);
    $("#psyFormLogin").val(curPsy.dec_login);
    $("#psyFormPas").val(curPsy.pas);
    $("#psyFormIdent").val(curPsy.dec_ident).prop("disabled", true);
    $("#psyFormCount").val(curPsy.count);

    $("input").prop("checked", false);
    curPsy.tests.forEach(function (testNumber) {
        $("#test"+testNumber).prop("checked", true)
    })

    $("#psyFormCheckDel").prop("checked", curPsy.pre_del !== ZeroDate);

}

function saveCurPsy() {
    curPsy.dec_login = $("#psyFormLogin").val();
    curPsy.pas = $("#psyFormPas").val();
    curPsy.count = $("#psyFormCount").val();
    curPsy.pre_del = $("#psyFormCheckDel").prop("checked");
    curPsy.tests = [];
    $('.testCheckbox:checked').each(function() { curPsy.tests.push(this.value); });
}



function validateFormData(login=$("#psyFormLogin"), pas=$("#psyFormPas"), ident=$("#psyFormIdent"), count=$("#psyFormCount")) {
    $("#psyFormBtnSave").prop("disabled", !(+validateText(login) + validatePas(pas) + validateText(ident) + validateNum(count) === 4));
}

function validateText(elem){
    let id =  elem.attr("id");
    if (id === "psyFormIdent" && curPsy) return true
    if(elem.val().match(/[^a-zA-Z0-9]/g) || !elem.val().length) {
        elem.toggleClass("is-invalid", true);
        $(`#${id}Msg`).text("Недопустимое значение");
        return false;
    }
    elem.toggleClass("is-invalid", false);
    return true;

}
function validatePas(elem){
    if(elem.val().match(/[^a-zA-Z0-9!"#$%&'()*,./:;=?@_`{|}~]/g) || elem.val().length < 9) {
        elem.toggleClass("is-invalid", true);
        $(`#${elem.attr("id")}Msg`).text("Недопустимый пароль. Он должен содержать не меннее 9 символов");
        return false;
    }
    elem.toggleClass("is-invalid", false);
    return true;

}
function validateNum(elem){
    if(elem.val().length && +elem.val() > 0) {
        elem.toggleClass("is-invalid", false);
        return true;

    }
    elem.toggleClass("is-invalid", true);
    $(`#${elem.attr("id")}Msg`).text("Неверное значение");
    return false;


}


function showPsy(key) {
    let psyTable = $("#psyTable");
    $('td').remove();
    if (!psyList) psyList = [];
    sort(psyList, key);

    fullCounter = { psyCount: psyList.length, whole: 0, not_yet: 0, clear: 0, danger: 0, msg: 0 };
    let grades, grade;

    for (let i = 0; i < psyList.length; i++) {
        let gradeCounter = { whole: 0, not_yet: 0, clear: 0, danger: 0, msg: 0 };
        grades = psyList[i].grades;
        for (let name in grades) {
            if (!grades.hasOwnProperty(name)) continue;
            grade = grades[name];
            gradeCounter.whole += grade.whole || 0;
            gradeCounter.not_yet += grade.not_yet || 0;
            gradeCounter.clear += grade.clear || 0;
            gradeCounter.danger += grade.danger || 0;
            gradeCounter.msg += grade.msg || 0;
        }
        gradeCounters[psyList[i].uid] = gradeCounter;

        let ownStats = `
            <span class="badge badge-light badge-pill" title="Количество испытуемых">${ gradeCounter.whole }</span>
            <span class="badge badge-secondary badge-pill" title="Еще не протестировано">${ gradeCounter.not_yet }</span>
            <span class="badge badge-success badge-pill" title="Вне групп риска">${ gradeCounter.clear }</span>
            <span class="badge badge-danger badge-pill" title="В группах риска">${ gradeCounter.danger }</span>`;
        if (gradeCounter.msg) {
            ownStats += `<span class="badge badge-warning badge-pill" title="Сообщения об удалении">${gradeCounter.msg}</span>`
        }
        if (!psyList[i].tests) psyList[i].tests = []
        let trPsy = $("<tr></tr>")
            .append($(`<td>${title(b64dec(psyList[i].login))}</td>`).click(function () { copyText(this) }))
            .append($(`<td>${psyList[i].pas}</td>`).click(function () { copyText(this) }))
            .append($(`<td>${psyList[i].count}</td>`))
            .append($(`<td>${ownStats}</td>`))
            .append($(`<td title="${psyList[i].tests.join(", ")}">${psyList[i].tests.length}</td>`))
            .append($(`<td>${ stamp2str(psyList[i].create_date) }</td>`))
            .append($(`<td><input type="button" class="btn btn-primary" onclick="showPsyInfoPage(${i})" value="Подробнее"></td>`));
        if (psyList[i].pre_del !== ZeroDate) trPsy.append($(`<td><i class="fa fa-trash" aria-hidden="true" title="Будет удален менее чем через 72 часа"></i></td>`));

        psyTable.append(trPsy);

        fullCounter.whole += gradeCounter.whole;
        fullCounter.not_yet += gradeCounter.not_yet;
        fullCounter.clear += gradeCounter.clear;
        fullCounter.danger += gradeCounter.danger;
        fullCounter.msg += gradeCounter.msg;

    }
    showStats(fullCounter);

}


function getPsyList(reloadTable= true, reloadMyData = false) {
    $("#loadingIcon").show();
    $.ajaxSetup({timeout:10000});
    $.get("/get_psy_list").done(function (response) {
        if (reloadMyData) getMyData()
        showMsg(response.msg, response.kind,function () {
            psyList = response.psyList;
            if (reloadTable) showPsy()
        });
    }).fail(function () { $("#loadingIcon").hide(); showMsg('Данные загрузить не удалось', "Err") });
}



function addNewPsy() {
    $.ajaxSetup({timeout:3000});
    $.post("/add_psy", $("#addPsyForm").serialize()).done(function (response) {
        showMsg(response.msg, response.kind,function () { clearPsyForm(); getPsyList(true); }, response.field);
    }).fail(function () {
        showMsg("Превышено время ожидания или произошла ошибка на стороне сервера! Операция не выполнена");
    })
}



function acceptDel(testeeUid, btn) {
    $.ajaxSetup({timeout:3000});
    $.post("/accept_del", {testeeUid: testeeUid}).done(function (response) {
        showMsg(response.msg, response.kind, function () {
            $(btn).toggleClass("invisAction", true);
            gradeCounters[curPsy.uid].msg -= 1;
            $("#stat_msg").text(gradeCounters[curPsy.uid].msg);
            if (gradeCounters[curPsy.uid].msg <= 0) $("#statsLinesMsg").toggleClass("d-flex", false).hide();
        });
    }).fail(function () {
        showMsg("Превышено время ожидания или произошла ошибка на стороне сервера! Операция не выполнена!");
    })
}


function editPsy() {
    let addPsyFormData = $('#addPsyForm').serializeArray();
    addPsyFormData.push({name: 'psyUid', value: curPsy.uid});
    $.ajaxSetup({timeout:3000});
    $.post("/edit_psy",  addPsyFormData).done(function (response) {
        showMsg(response.msg, response.kind,function () {
            $("#psyFormBtnSave").prop("disabled", true);
            needToReload = true;
            saveCurPsy();
        }, response.field);
    }).fail(function () {
        showMsg("Превышено время ожидания или произошла ошибка на стороне сервера! Операция не выполнена");
    })
}



function showGrades(key) {
    let gradeTable = $("#gradeTable");
    $('#gradeTable td').remove();
    sort(gradeList, key);

    let gradeCounter = { gradeCount: gradeList.length, whole: 0, not_yet: 0, clear: 0, danger: 0, msg: 0 };

    for (let i = 0; i < gradeList.length; i++) {
        let grade = gradeList[i];
        gradeCounter.whole += grade.whole;
        gradeCounter.not_yet += grade.not_yet;
        gradeCounter.clear += grade.clear;
        gradeCounter.danger += grade.danger;
        gradeCounter.msg += grade.msg;

        let trGrade = $("<tr></tr>").append($(`<td>${grade.dec_name}</td>`))
            .append($("<td></td>").append($(`<span class="badge badge-Light badge-pill">${grade.whole}</span>`)))
            .append($("<td></td>").append($(`<span class="badge badge-secondary badge-pill">${grade.not_yet}</span>`)))
            .append($("<td></td>").append($(`<span class="badge badge-success badge-pill">${grade.clear}</span>`)))
            .append($("<td></td>").append($(`<span class="badge badge-danger badge-pill">${grade.danger}</span>`)))
            .append($(`<td><input type="button" class="btn btn-primary" onclick="showGradePage(${i})" value="Просмотреть"></td>`));
        if (grade.msg) {
            trGrade.append(`<td><span class="btn btn-warning my-2 my-sm-0" style="cursor: default" title="В этом классе есть запросы на удаление результата">
                <i class="fa fa-exclamation-triangle" aria-hidden="true"></i>&nbsp;${grade.msg}</span></td>`);
        }
        gradeTable.append(trGrade);
    }
    showStats(gradeCounter);
    gradeCounters[curPsy.uid] = gradeCounter;
}



function getGradeList(reloadTable = true) {
    $("#loadingIcon").show();
    $.ajaxSetup({timeout:10000});
    $.get("/get_user_data", { psyUid: curPsy.uid}).done(function (response) {
        showMsg(response.msg, response.kind,function () {
            gradeList = [];
            if (response.userData && $("#psyFormBtnSave").prop("disabled")) { curPsy = response.userData; setToDefault()}
            if (!response.userData.grades) { showGrades(); return  }
            for (let name in response.userData.grades) {
                if (!response.userData.grades.hasOwnProperty(name)) continue;
                gradeList.push({
                    name: name,
                    dec_name: b64dec(name).toUpperCase(),
                    whole: response.userData.grades[name].whole || 0,
                    not_yet: response.userData.grades[name].not_yet || 0,
                    clear: response.userData.grades[name].clear || 0,
                    danger: response.userData.grades[name].danger || 0,
                    msg: response.userData.grades[name].msg || 0})
            }
            if (reloadTable) showGrades()
        });
    }).fail(function () { $("#loadingIcon").hide(); showMsg('Данные загрузить не удалось', "Err")


    });
}




function showTestees(key="result", reverseResults=true) {
    let testeeTable = $("#testeeTable");
    $('#testeeTable td').remove();
    if (!testeeList) testeeList =[];
    sort(testeeList, key, true);

    let gradeCounter = { whole: testeeList.length, not_yet: 0, clear: 0, danger: 0, msg: 0 };

    for (let i = 0; i < testeeList.length; i++) {
        let testee = testeeList[i];
        gradeCounter[resultDecode[testee.result][0]] += 1;

        let trTestee = $("<tr></tr>")
            .append($(`<td><span class="badge badge-${resultDecode[testee.result][1]} badge-pill">${resultDecode[testee.result][2]}</span></td>`))
            .append($(`<td>${b64dec(testee.login)}</td>`).click(function () { copyText(this) }))
            .append($(`<td>${testee.pas}</td>`).click(function () { copyText(this) }))
            .append($(`<td>${stamp2str(testee.create_date)}</td>`));

        if (testee.msg) {
            gradeCounter.msg += 1;
            trTestee.append(`<td class="actionTd">
                <div  id="delBtn${i}">
                    <span data-toggle="dropdown" aria-haspopup="true" aria-expanded="false" class="btn btn-outline-warning my-2 my-sm-0" style="cursor: pointer;"
                        title="Нажмите, для просмотра сообщения об удалении"><i class="fa fa-exclamation-triangle" aria-hidden="true"></i>
                    </span>
                    <div class="dropdown-menu">
                        <div class="card border-0 shadow" id="show_msg_card">
                            <div class="card-body">
                                <h5 class="card-title">Причина удаления</h5>
                                <textarea readonly style="width: 600px" class="form-control" rows="2" required maxlength="500" aria-describedby="stopLen" name="reason">${b64dec(testee.msg)}</textarea>
                                <br>
                                <input type="button" class="btn btn-primary mr-1" value="Подтвердить удаление" onclick="acceptDel('${testee.uid}', delBtn${i})">
                            </div>
                        </div>
                    </div>
               </div></td>`);
        } else {
            trTestee.append(`<td class="actionTd">&nbsp;</td>`);
        }
        testeeTable.append(trTestee);
    }
    showStats(gradeCounter);
    gradeCounters[curPsy.uid] = gradeCounter;

}


function getTesteeList(reloadTable= true) {
    $("#loadingIcon").show();
    $.ajaxSetup({timeout:10000});
    $.get("/get_testee_list", {psyUid: curPsy.uid, grade: curGrade.dec_name}).done(function (response) {
        testeeList = response.testeeList
        if (reloadTable) showTestees()
    }).fail(function () { $("#loadingIcon").hide(); showMsg('Данные загрузить не удалось', "Err") });
}



function clearPsyForm() {
    $("#psyFormBtnSave").prop("disabled", true);
    $("#psyFormLogin").val("");
    $("#psyFormPas").val(generatePas(12));
    $("#psyFormIdent").val("");
    $("#psyFormCount").val("");
    $(".testCheckbox:checked").prop("checked", false)
}


function showPsyInfoPage(psyIdx) {
    if (curGrade) { curGrade = undefined; $("#add_psy_card").slideDown(); $("#testeeTablePlace").hide(); }
    else curPsy = psyList[psyIdx];
    let btnSave = $("#psyFormBtnSave");
    setToDefault();
    $("#psyTablePlace").hide();
    $("#statsLinesPsyCount").removeClass("d-flex").hide();

    $("#gradeTablePlace").show();
    $("#psyFormBtnDef").show();
    $("#psyFormPlaceDel").show();
    $("#barBtnBack").off("click").click(showAdminMainPage).show();

    $("#psyFormTitle").text("Редактировать Психолога");
    $("#statsCardTitle").text(`${curPsy.dec_login} | Статистика`);
    btnSave.off("click").click(function () { rareCall(editPsy) }).val("Сохранить");
    $("#statsCardBtnRefresh").off("click").click(function () { rareCall(getGradeList) });

    setDownloadLinks("", curPsy.uid);

    $("input").toggleClass("is-invalid", false);
    btnSave.prop("disabled", true);
    showStats(gradeCounters[curPsy.uid]);
    getGradeList();

}


function showAdminMainPage() {
    clearPsyForm();
    curPsy = undefined;
    let saveBtn = $("#psyFormBtnSave");
    $("#psyTablePlace").show();
    $("#statsLinesPsyCount").addClass("d-flex").show();

    $("#gradeTablePlace").hide();
    $("#psyFormBtnDef").hide();
    $("#psyFormPlaceDel").hide();
    $("#barBtnBack").hide();

    $("#psyFormTitle").text("Добавить психолога");
    $("#statsCardTitle").text(`Полная статистика`);

    saveBtn.off("click").click(function () { rareCall(addNewPsy) }).val("Добавить психолога");
    $("#statsCardBtnRefresh").off("click").click(function () { rareCall(getPsyList) });

    $("input").toggleClass("is-invalid", false);
    saveBtn.prop("disabled", true);
    $("#psyFormIdent").val("").prop("disabled", false);
        setDownloadLinks();

    showStats(fullCounter);
    if (needToReload) { getPsyList(); needToReload = false}

}


function showGradePage(gradeIdx) {
    curGrade = gradeList[gradeIdx];
    $("#add_psy_card").slideUp();

    $("#gradeTablePlace").hide();

    $("#testeeTablePlace").show();
        $("#barBtnBack").off("click").click(showPsyInfoPage);

    $("#gradeName").text(curGrade.dec_name);

    $("#statsCardTitle").text(`${curGrade.dec_name} | Статистика`);
    $("#statsCardBtnRefresh").off("click").click(function () { rareCall(getTesteeList) });
    setDownloadLinks(curGrade.dec_name, curPsy.uid);
    showStats(curGrade);
    getTesteeList();
    needToReload = true;
}


function getMyData() {
    $.ajaxSetup({timeout:2000});
    $.get("/get_user_data", {isMy:"true"}).done(function (response) {
        showMsg(response.msg, response.kind, function () {
            myData = response.userData;
            myData.dec_login = title(b64dec(myData.login));
            setLogin(myData);
        })}).fail(function () {
        showMsg('Данные текущего пользователя загрузить не удалось', "Fatal")
    });
    }


$("#psyTablePlace").ready(function () { getPsyList(true, true) });
$("#psyFormBtnSave").ready(function () { $("#psyFormBtnSave").click(addNewPsy) });
$("#statsCardBtnRefresh").ready(function () { $("#statsCardBtnRefresh").click(function () { rareCall(getPsyList) }) });
$("#statsCardBtnDownload").ready(function () { setDownloadLinks() });




