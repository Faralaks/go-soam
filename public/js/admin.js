let testeeList;
let myData;



function delResult(testeeUid) {
    $.ajaxSetup({timeout:3000});
    $.post("/accept_del", {testeeUid: testeeUid}).done(function (response) {
        showMsg(response.msg, response.kind, function () {
        });
    }).fail(function () {
        showMsg("Превышено время ожидания или произошла ошибка на стороне сервера! Операция не выполнена!");
    })
}


function getTesteeList(reloadTable= true) {
    //$("#loadingIcon").show();
    $.ajaxSetup({timeout:10000});
    $.get("/get_testee_list").done(function (response) {
        testeeList = response.testeeList;
        if (reloadTable) showTestees()
    }).fail(function () { $("#loadingIcon").hide(); showMsg('Данные загрузить не удалось', "Err") });
}



function showTestees(key="result", reverseResults=true) {
    let testeeTable = $("#testeeTable");
    $('#testeeTable td').remove();
    if (!testeeList) testeeList =[];
    sort(testeeList, key, true);
    let td = `<td></td>`;
    for (let i = 0; i < testeeList.length; i++) {
        let testee = testeeList[i];
        let trTestee = $(`<tr class="greyRow" onclick="testeePage(${i}, $(this))"></tr>`);
        trTestee.append($(td).text(b64dec(testee.login)).click(function () { copyText(this) }));
        trTestee.append($(td).text(b64dec(testee.name)).click(function () { copyText(this) }));
        trTestee.append($(td).text(testee.ege));
        trTestee.append($(td).text(testee.grade));
        trTestee.append($(td).text(stamp2str(testee.create_date)));

        testeeTable.append(trTestee);

    }

}



function testeePage(testeeIdx, row) {
    let testee = testeeList[testeeIdx];
    if (testee.dataRow) {
        testee.dataRow.slideToggle();
        return;
    }

    let newTr = $(`<tr></tr>`);
    let newTd = $(`<td style="" colspan="5" class="dn"></td>`);
    let newDiv =  $(`<div style="height: 100%; width: 100%" class="container card"></div>`);
    console.log(testee.result)

    for (let blank in testee.result) {
        if (!testee.result.hasOwnProperty(blank)) continue;
        console.log("blank", blank)
        visFuncs[blank](testee.result[blank], newDiv)
    }


    newTd.append(newDiv);
    newTr.append(newTd);
    row.after(newTr);
    testee.dataRow = row.next().children();
    testee.dataRow.slideDown();



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




