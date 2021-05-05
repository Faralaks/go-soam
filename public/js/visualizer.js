

let visFuncs = {
    BPAQ: visBPAQ,
    ITO: visITO,
};



function getColorForLF(num) {
    if (num > 5) return "danger"
    return "none"
}

let dotCords = {
    dc: {x: 400, y: 400},
    d0: {x: 270*Math.cos(0)+400, y: 400-270*Math.sin(0)},
    d45: {x: 270*Math.cos(Math.PI/4)+400, y: 400-270*Math.sin(Math.PI/4)},
    d90: {x: 270*Math.cos(Math.PI/2)+400, y: 400-270*Math.sin(Math.PI/2)},
    d135: {x: 270*Math.cos(3*Math.PI/4)+400, y: 400-270*Math.sin(3*Math.PI/4)},
    d180: {x: 270*Math.cos(Math.PI)+400, y: 400-270*Math.sin(Math.PI)},
    d225: {x: 270*Math.cos(5*Math.PI/4)+400, y: 400-270*Math.sin(5*Math.PI/4)},
    d270: {x: 270*Math.cos(3*Math.PI/2)+400, y: 400-270*Math.sin(3*Math.PI/2)},
    d315: {x: 270*Math.cos(7*Math.PI/4)+400, y: 400-270*Math.sin(7*Math.PI/4)},
};

function divideLine(s, num) {
    num = +num;
    let f = dotCords.dc;
    let l = num/(9-num);
    let x = (f.x + l*s.x)/(1+l);
    let y = (f.y + l*s.y)/(1+l);
    return {x: x, y: y}
    //return `cx='${x}'  cy='${y}'`
}



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


function visITO (data, container) {
    console.log(data);

    let resDotCords = {
        i: divideLine(dotCords.d0, +data.i),
        ii: divideLine(dotCords.d315, +data.ii),
        iii: divideLine(dotCords.d270, +data.iii),
        iv: divideLine(dotCords.d225, +data.iv),
        v: divideLine(dotCords.d180, +data.v),
        vi: divideLine(dotCords.d135, +data.vi),
        vii: divideLine(dotCords.d90, +data.vii),
        viii: divideLine(dotCords.d45, +data.viii),
    };

    let newTbl = $(`
        <table class="table text-center"><tr><td colspan="3" ><b>Индивидуально-типологический опросник</b></td></tr> 
        <tr>
            <td colspan="2">Ложь: <span class="badge badge-pill badge-${getColorForLF(+data.L)}">${data.L}</span></td>
            <td colspan="2">Аггравация: <span class="badge badge-pill badge-${getColorForLF(+data.F)}">${data.F}</span></td>
            </tr>
            <tr>
                <td>Экстраверсия: ${data.i}</td>
                <td>Спонтанность: ${data.ii}</td>
                <td>Агрессивность: ${data.iii}</td>
                <td>Ригидность: ${data.iv}</td>
        </tr>
            <tr>
                <td>Интроверсия: ${data.v}</td>
                <td>Сензитивность: ${data.vi}</td>
                <td>Тревожность: ${data.vii}</td>
                <td>Эмотивность: ${data.viii}</td>
        </tr>
        <tr>
            <td colspan="4" ><svg height="800" width="800" xmlns="http://www.w3.org/2000/svg">
                <g>
                    <rect id="back" height="800" width="800" y="0" x="0" stroke-width="1.5" stroke="#000" fill="#fff"/>
                    <circle id="range9" cx="400" cy="400" r="270" stroke-width="1.5" stroke="#000" fill="#fff"/>
                    <circle id="range7" cx="400" cy="400" r="210" stroke-width="1.5" stroke="#000" fill="#fff"/>
                    <circle id="range4" cx="400" cy="400" r="120" stroke-width="1.5" stroke="#000" fill="#fff"/>
           
                     <path d="M ${dotCords.d0.x} ${dotCords.d0.y}  L ${dotCords.d180.x} ${dotCords.d180.y}" fill="transparent" stroke="black"/>
                     <path d="M ${dotCords.d45.x} ${dotCords.d45.y}  L ${dotCords.d225.x} ${dotCords.d225.y}" fill="transparent" stroke="black"/>
                     <path d="M ${dotCords.d90.x} ${dotCords.d90.y}  L ${dotCords.d270.x} ${dotCords.d270.y}" fill="transparent" stroke="black"/>
                     <path d="M ${dotCords.d135.x} ${dotCords.d135.y}  L ${dotCords.d315.x} ${dotCords.d315.y}" fill="transparent" stroke="black"/>
                     
                     <circle id="d0" cx="${dotCords.d0.x}" cy="${dotCords.d0.y}" r="5" fill="#000"/>
                     <circle id="d45" cx="${dotCords.d45.x}" cy="${dotCords.d45.y}" r="5" fill="#000"/>
                     <circle id="d90" cx="${dotCords.d90.x}" cy="${dotCords.d90.y}" r="5" fill="#000"/>
                     <circle id="d135" cx="${dotCords.d135.x}" cy="${dotCords.d135.y}" r="5" fill="#000"/>
                     <circle id="d180" cx="${dotCords.d180.x}" cy="${dotCords.d180.y}" r="5" fill="#000"/>
                     <circle id="d2250" cx="${dotCords.d225.x}" cy="${dotCords.d225.y}" r="5" fill="#000"/>
                     <circle id="d270" cx="${dotCords.d270.x}" cy="${dotCords.d270.y}" r="5" fill="#000"/>
                     <circle id="d315" cx="${dotCords.d315.x}" cy="${dotCords.d315.y}" r="5" fill="#000"/>
                     
                     <circle id="center" cx="400" cy="400" r="5" fill="#000"/>
                     
                     <circle id="i" cx="${resDotCords.i.x}" cy="${resDotCords.i.y}" r="5" fill="#f00"/>
                     <circle id="ii" cx="${resDotCords.ii.x}" cy="${resDotCords.ii.y}" r="5" fill="#f00"/>
                     <circle id="iii" cx="${resDotCords.iii.x}" cy="${resDotCords.iii.y}" r="5" fill="#f00"/>
                     <circle id="iv" cx="${resDotCords.iv.x}" cy="${resDotCords.iv.y}" r="5" fill="#f00"/>
                     <circle id="v" cx="${resDotCords.v.x}" cy="${resDotCords.v.y}" r="5" fill="#f00"/>
                     <circle id="vi" cx="${resDotCords.vi.x}" cy="${resDotCords.vi.y}" r="5" fill="#f00"/>
                     <circle id="vii" cx="${resDotCords.vii.x}" cy="${resDotCords.vii.y}" r="5" fill="#f00"/>
                     <circle id="viii" cx="${resDotCords.viii.x}" cy="${resDotCords.viii.y}" r="5" fill="#f00"/>
                     
                     <polygon points="${resDotCords.i.x},${resDotCords.i.y} ${resDotCords.ii.x},${resDotCords.ii.y} ${resDotCords.iii.x},${resDotCords.iii.y}  ${resDotCords.iv.x},${resDotCords.iv.y} 
                     ${resDotCords.v.x},${resDotCords.v.y} ${resDotCords.vi.x},${resDotCords.vi.y} ${resDotCords.vii.x},${resDotCords.vii.y}  ${resDotCords.viii.x},${resDotCords.viii.y}"
                     style="fill:rgba(97,168,255,0.4);stroke:black;stroke-width:1;fill-rule:evenodd;" />
                     
                     
                     
                     

                </g>
            
            </svg></td>
        </tr>
        
        
        </table>`);

    container.append(newTbl)
}


