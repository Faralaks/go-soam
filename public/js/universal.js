let lastKey;
let nextTime = 0;

let fieldNamesDecode = {
    Login: "Логин",
    Ident: "Идентификатор"
};
let resultDecode = [];
resultDecode[-1] = ["not_yet", "secondary", "Нет результата"];
resultDecode[0] = ["clear", "success", "Вне группы"];
resultDecode[1] = ["danger", "danger", "В группе"];

const NotYetResult = -1
const ClearResult = 0
const DangerResult = 1

var ZeroDate = "0001-01-01T00:00:00Z"


let shuffledAlf = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '!',  ';', '%', ':',  '@', '(', ')', '_', '=']
                                .sort(function(){ return Math.random() - 0.5; }).join("");


function showMsg(msg, kind,  sucFunc=function () {}, field="") {
    switch (kind) {
        case "Suc":
            $("#formBtn").prop("disabled", true);
            $("#formSignSuc").fadeIn(1000).delay(3000).fadeOut(500);
            sucFunc();
            return;
        case "DuplicatedField":
            let capitalizedField = field.replace(/(^|\s)\S/g, l => l.toUpperCase());
            $(`#psyForm${capitalizedField}`).toggleClass("is-invalid", true);
            $(`#psyForm${capitalizedField}Msg`).text(`Такой ${fieldNamesDecode[capitalizedField]} уже существует`);
            return;
        case "Fatal":
            $("#fatalMsg").text("Ошибка: "+msg).fadeIn(700).delay(6000).fadeOut(4000);
            return;
        case "BadUpdate":
            $("#updateErr").show();
            return;
        case "Relogin":
            alert("Ваша сессия истекла, необходимо войти повторно.")
            document.location.replace("/");
            return;

        default:
            sucFunc(); return;
    }

}


function sort(list, key, reverseResults=false) {
    if (key) {
        if (key === lastKey) { reverse *= -1; }
        else { reverse = 1; lastKey = key; }
        if (reverseResults) {reverse = -1}

        list.sort(function (a, b) {
            if (a[key] > b[key]) { return reverse; }
            if (a[key] < b[key]) { return -1*reverse; }
            return 0;
        });
    }
    return list;

}



function generatePas(len){
    let pas = "";
    for (let i=0; i<len; i++){
        pas += shuffledAlf.charAt(Math.floor(Math.random() * shuffledAlf.length));
    }
    return pas;
}


function stamp2str(timestamp){
  let date = new Date(timestamp * 1000);
  return date.toLocaleString().replace(", ", "<br>");
  l
}

function b64enc(text) {
	return window.btoa(unescape(encodeURIComponent(text)));
}

function b64dec(text) {
	return decodeURIComponent(escape(window.atob(text)));
}

function title(str) {
    return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase()
}

function copyText(el) {
    var $tmp = $("<textarea>");
    $("body").append($tmp);
    $tmp.val($(el).text()).select();
    document.execCommand("copy");
    $tmp.remove();
}

function rareCall(func, delay=1500) {
    if (nextTime < Date.now()) {
        func();
        nextTime = Date.now() + delay;
    }
}


function setLogin(user) {
    $("#barPlace").fadeIn(300);
    $("#loginPlace").text(user.dec_login);
    if (user.modifiedDate && Math.ceil(((Date.now() / 1000 | 0) - user.modifiedDate) / 60) > 50) $("#pasWarning").show(); //86400
}

function setDownloadLinks(grade="", owner="") {
    if (grade === "") { $("#downloadDocx").hide(); }
    else $("#downloadDocx").attr("href", `/download?psyUid=${encodeURIComponent(owner)}&grade=${encodeURIComponent(grade)}&target=${encodeURIComponent('not_yet')}`).show();

    $("#downloadXlsx").attr("href", `/download?psyUid=${encodeURIComponent(owner)}&grade=${encodeURIComponent(grade)}&target=${encodeURIComponent('done')}`);
}

function validateNewPas(elem){
    if(elem.val().match(/[^a-zA-Z0-9!"#$%&'()*,./:;=?@_`{|}~]/g) || elem.val().length < 9) {
        elem.toggleClass("is-invalid", true);
        $(`#${elem.attr("id")}Msg`).show();
        return false;
    }
    elem.toggleClass("is-invalid", false);
    return true;

}

function confPas(pas, conf) {
    if (pas.val() !== conf.val()) {
        conf.toggleClass("is-invalid", true);
        $(`#${conf.attr("id")}Msg`).show();
        return false;
    } else {
        conf.toggleClass("is-invalid", false);
        $(`#${conf.attr("id")}Msg`).hide();

        return true;
    }
}

function sendNewData() {
    $.ajaxSetup({timeout:2000});
    $.post("/edit_user_data", $("#newPasForm").serialize()).done(function (response) {
        showMsg(response.msg, response.kind)}).fail(function () {
        showMsg('Неизвестная ошибка во время запроса, возможно, соединение с сервером потеряно.', "Fatal")
    });
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

function validateNum(elem){
    if(elem.val().length && +elem.val() > 0) {
        elem.toggleClass("is-invalid", false);
        return true;

    }
    elem.toggleClass("is-invalid", true);
    $(`#${elem.attr("id")}Msg`).text("Неверное значение");
    return false;

}