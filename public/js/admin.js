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

function showTestees(key="result", reverseResults=true) {
    let testeeTable = $("#testeeTable");
    $('#testeeTable td').remove();
    if (!testeeList) testeeList =[];
    sort(testeeList, key, true);

    for (let i = 0; i < testeeList.length; i++) {
        let testee = testeeList[i];
        let trTestee = $("<tr></tr>")
            .append($(`<td>${b64dec(testee.login)}</td>`).click(function () { copyText(this) }))
            .append($(`<td>${testee.pas}</td>`).click(function () { copyText(this) }))
            .append($(`<td>${testee.ege}</td>`))
            .append($(`<td>${testee.grade}</td>`))
            .append($(`<td>${stamp2str(testee.create_date)}</td>`));

        testeeTable.append(trTestee);
    }

}


function getTesteeList(reloadTable= true) {
    //$("#loadingIcon").show();
    $.ajaxSetup({timeout:10000});
    $.get("/get_testee_list").done(function (response) {
        testeeList = response.testeeList
        if (reloadTable) showTestees()
    }).fail(function () { $("#loadingIcon").hide(); showMsg('Данные загрузить не удалось', "Err") });
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


$("#testeeTablePlace").ready(function () { getMyData(); getTesteeList(true) });




