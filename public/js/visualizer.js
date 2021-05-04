

let visFuncs = {
    BPAQ: visBPAQ
};

function visBPAQ (data, container) {
    console.log(data)

    let newTbl = $(`
        <table class="table text-center"><tr><td colspan="3" ><b>Личностный опросник агрессивности Басса-Перри</b></td></tr> <tr>
            <td>Физическая агрессия: ${data.aggression}</td>
            <td>Гнев: ${data.anger}</td>
            <td>Враждебность: ${data.hostility}</td>
        </tr></table>`);

    container.append(newTbl)
}