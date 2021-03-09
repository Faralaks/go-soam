let curPsy, curGrade;
let gradeList, testeeList, nameList;
let psyCounter, gradeCounter;


function validateFormData() {
    $("#addFormBtnAdd").prop("disabled", !(+validateText($("#addFormName")) + validateNum($("#addFormCount")) === 2));
}

function validateText(elem){
    if(elem.val().match(/[^a-zA-Zа-яА-Я0-9«»“”"„ ]/g) || !elem.val().length) {
        elem.toggleClass("is-invalid", true);
        $(`#${elem.attr("id")}Msg`).text("Недопустимое значение");
        return false;
    }
    elem.toggleClass("is-invalid", false);
    return true;

}
function validateNum(elem){
    if(elem.val().length && +elem.val() > 0 && +elem.val() <= curPsy.count) {
        elem.toggleClass("is-invalid", false);
        return true;
    }
    elem.toggleClass("is-invalid", true);
    $(`#${elem.attr("id")}Msg`).text("Недопустимое значение");
    return false;


}




function addHelperLines() {
    let value = $(this).val().toLowerCase();
    $("#helperList").empty();
    for (let i = 0; i < nameList.length; i++) {
        if (i > 3) break;
        if (nameList[i].toLowerCase().indexOf(value) === -1) continue;
        $("#helperList").append($(`<span  class="list-group-item list-group-item-action list-group-item-light" onclick="if (!curGrade) addFormName.value = this.textContent">${nameList[i]}</span>`))
    }
}


function renderGradeList(list) {
    gradeList = [];
    nameList = [];
    if (!list) return
    for (let name in list) {
        if (!list.hasOwnProperty(name)) continue;
        nameList.push(b64dec(name));
        gradeList.push({name: name,
        dec_name: b64dec(name),
        whole: list[name].whole || 0,
        not_yet: list[name].not_yet || 0,
        clear: list[name].clear || 0,
        danger: list[name].danger || 0,
        msg: list[name].msg || 0})
    }
    nameList.sort();
}

function clearTesteeForm() {
    $("#addFormBtnAdd").prop("disabled", true);
    if (!curGrade) $("#addFormName").prop("readonly", false).val("");
    else $("#addFormName").prop("readonly", true).val(curGrade.dec_name);
    $("#addFormCount").val("");
    $("input").toggleClass("is-invalid", false);
}


function showGrades(key) {
    let gradeTable = $("#gradeTable");
    if (!gradeList) return
    $('#gradeTable td').remove();
    sort(gradeList, key);

    psyCounter = { gradeCount: gradeList.length, whole: 0, not_yet: 0, clear: 0, danger: 0};

    for (let i = 0; i < gradeList.length; i++) {
        let grade = gradeList[i];
        psyCounter.whole += grade.whole;
        psyCounter.not_yet += grade.not_yet;
        psyCounter.clear += grade.clear;
        psyCounter.danger += grade.danger;

        let trGrade = $("<tr></tr>").append($(`<td>${grade.dec_name}</td>`))
            .append($("<td></td>").append($(`<span class="badge badge-Light badge-pill">${grade.whole}</span>`)))
            .append($("<td></td>").append($(`<span class="badge badge-secondary badge-pill">${grade.not_yet}</span>`)))
            .append($("<td></td>").append($(`<span class="badge badge-success badge-pill">${grade.clear}</span>`)))
            .append($("<td></td>").append($(`<span class="badge badge-danger badge-pill">${grade.danger}</span>`)))
            .append($(`<td><input type="button" class="btn btn-success" onclick="showGradePage(${i})" value="Просмотреть"></td>`));
        gradeTable.append(trGrade);
    }
    showStats(psyCounter);
}


function deleteResult(testeeIdx, btn) {
    let testee = testeeList[testeeIdx];
    let delMsg = $("#delReasonField"+testeeIdx).val();
    $.ajaxSetup({timeout:2000});
    $.post("/del_result", {testeeUid: testee.uid, msg: delMsg, grade: curGrade.dec_name}).done(function (response) {
        showMsg(response.msg, response.kind, function () {
            $(btn).hide();
            $(btn).parent().append($(`<span  title="Этот результат был удален по вашей причине: ${delMsg}"><i class="fa fa-trash" aria-hidden="true"></i></span>`))
            gradeCounter[resultDecode[testee.result][0]] -= 1;
            $("#stat_"+resultDecode[testee.result][0]).text(gradeCounter[resultDecode[testee.result][0]]);
            gradeCounter.not_yet += 1;
            $("#stat_not_yet").text(gradeCounter.not_yet);
            $("#resultPlace"+testeeIdx).text(resultDecode[NotYetResult][2]).toggleClass("badge-"+resultDecode[testee.result][1], false).toggleClass("badge-secondary", true);
            testeeList[testeeIdx].result = NotYetResult;
            testeeList[testeeIdx].msg = delMsg;
        });
    }).fail(function () { $("#loadingIcon").hide(); showMsg('Данные загрузить не удалось', "Fatal") });
}


function showTestees(key="result", reverseResults= true) {
    let testeeTable = $("#testeeTable");
    $('#testeeTable td').remove();
    if (!testeeList ) return
    sort(testeeList, key, reverseResults);

    gradeCounter = { whole: testeeList.length, not_yet: 0, clear: 0, danger: 0 };
    for (let i = 0; i < testeeList.length; i++) {
        let testee = testeeList[i];
        testee.dec_login = b64dec(testee.login)
        gradeCounter[resultDecode[testee.result][0]] += 1;
        let trTestee = $(`<tr onmouseenter="let btn = $(\'#delBtn${i}\'); if (btn) {btn.css('visibility', 'visible')}" onmouseleave="let btn = $(\'#delBtn${i}\'); if (btn) {btn.css('visibility', 'collapse')}"></tr>`)
            .append($(`<td><span id="resultPlace${i}" class="badge badge-${resultDecode[testee.result][1]} badge-pill">${resultDecode[testee.result][2]}</span></td>`))
            .append($(`<td>${b64dec(testee.login)}</td>`).click(function () { copyText(this) }))
            .append($(`<td>${testee.pas}</td>`).click(function () { copyText(this) }))
            .append($(`<td>${stamp2str(testee.create_date)}</td>`));

        if (testee.result !== NotYetResult ) {
            trTestee.append(`<td class="actionTd">
                <span class="btn btn-outline-danger my-2 my-sm-0" id="delBtn${i}" data-toggle="modal" data-target="#delResultModel${i}"  style="visibility: collapse"
                      title="Нажмите, чтобы удалить результат"> <i class="fa fa-trash" aria-hidden="true"></i>
                </span>
                <div class="modal fade" id="delResultModel${i}" tabindex="-1" role="dialog" aria-labelledby="delResultModelTitle${i}" aria-hidden="true">
                    <div class="modal-dialog modal-dialog-centered" role="document" style="text-align: left">
                        <div class="modal-content">
                            <div class="modal-header">
                                <h5 class="model-title" id="delResultModelTitle${i}">Введите причину удаления результата испытуемого ${testee.dec_login} из ${curGrade.dec_name}</h5>
                                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                                    <span aria-hidden="true" id="modelExitBtn${i}">&times;</span>
                                </button>
                            </div>
                            <div class="modal-body">
                                <textarea id="delReasonField${i}" class="form-control" rows="4" required maxlength="500" placeholder="Причина удаления" aria-describedby="stopLen" name="reason"></textarea>
                                <small id="stopLen" class="form-text text-muted">Не более 500 символов</small>
                                <br>
                                <input type="button" class="btn btn-danger pull-right" value="Удалить" onclick="$(modelExitBtn${i}).trigger('click'); deleteResult('${i}', delBtn${i})">
                            </div>
                        </div>
                    </div>
                </div></td>`);
        }
        if (testee.msg) {
            trTestee.append(`<td>
                <span  title="Этот результат был удален по вашей причине: ${b64dec(testee.msg)}">
                    <i class="fa fa-trash" aria-hidden="true"></i>
                </span>
            </td>`)
        } else {
            trTestee.append(`<td class="actionTd">&nbsp;</td>`);
        }
        testeeTable.append(trTestee);
    }
    showStats(gradeCounter);

}


function getTesteeList(reloadTable= true) {
    $("#loadingIcon").show();
    $.ajaxSetup({timeout:10000});
    $.get("/get_testee_list", {psyUid: curPsy.uid, grade: curGrade.dec_name}).done(function (response) {
        showMsg(response.msg, response.kind, function () {
            testeeList = response.testeeList
            if (reloadTable) showTestees() });
    }).fail(function () { $("#loadingIcon").hide(); showMsg('Данные загрузить не удалось', "Fatal") });
}



function getUserData() {
    $("#loadingIcon").show();
    $.ajaxSetup({timeout:2000});
    $.get("/get_user_data").done(function (response) {
        showMsg(response.msg, response.kind, function () {
            curPsy = response.userData;
            curPsy.dec_login = title(b64dec(curPsy.login))
            setLogin(curPsy);
            $("#countPlace").text(curPsy.count);
            $("#addFormCount").prop("max", curPsy.count);
            renderGradeList(curPsy.grades);
            showGrades() });
    }).fail(function () { $("#loadingIcon").hide(); showMsg('Данные загрузить не удалось', "Fatal") });
}


function addTestees() {
    $.ajaxSetup({timeout:3000});
    $.post("/add_testees", $("#addTesteesForm").serialize()).done(function (response) {
        showMsg(response.msg, response.kind, function () {
            if (!curGrade) { renderGradeList(response.newGrades); showGrades() }
            let formCount = $("#addFormCount");
            curPsy.count -= +formCount.val();
            $("#countPlace").text(curPsy.count);
            formCount.prop("max", curPsy.count);
            clearTesteeForm();
            if (curGrade) { getTesteeList() } });
    }).fail(function () {
        showMsg("Превышено время ожидания или произошла ошибка на стороне сервера! Операция не выполнена", "Fatal");
    })
}



function showGradePage(gradeIdx) {
    curGrade = gradeList[gradeIdx];
    clearTesteeForm();
    $("#gradeTablePlace").hide();

    $("#testeeTablePlace").show();
    $("#barBtnBack").off("click").click(showPsyMainPage).show();

    $("#gradeName").text(curGrade.dec_name);

    $("#statsCardTitle").text(`${curGrade.dec_name} | Статистика`);
    $("#statsCardBtnRefresh").off("click").click(function () { rareCall(getTesteeList) });
    setDownloadLinks(curGrade.dec_name);
    showStats(curGrade);
    getTesteeList();
}


function showPsyMainPage() {
    curGrade = undefined;
    testeeList = undefined;
    gradeCounter = undefined;
    clearTesteeForm();

    $("#gradeTablePlace").show();

    $("#testeeTablePlace").hide();
    $("#barBtnBack").off("click").hide();

    $("#statsCardTitle").text("Общая статистика");
    $("#statsCardBtnRefresh").off("click").click(function () { rareCall(getUserData) });
    setDownloadLinks();
    showStats(psyCounter);
    getUserData();
}





$("#gradeTablePlace").ready(getUserData);
$("#statsCardBtnRefresh").ready(function () { $("#statsCardBtnRefresh").click(function () { rareCall(getUserData) }) });
$("#helperList").ready(function() { $("#addFormName").on("input", addHelperLines); });
$("#statsCardBtnDownload").ready(function () { setDownloadLinks() });
