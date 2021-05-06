

let visFuncs = {
    BPAQ: visBPAQ,
    ITO: visITO,
};



function getColorForLF(num) {
    if (num > 5) return "danger"
    return "none"
}

let dotCords = {
    dc: {x: 400, y: 330},
    i: {x: 270*Math.cos(0)+400, y: 330-270*Math.sin(0)},
    viii: {x: 270*Math.cos(Math.PI/4)+400, y: 330-270*Math.sin(Math.PI/4)},
    vii: {x: 270*Math.cos(Math.PI/2)+400, y: 330-270*Math.sin(Math.PI/2)},
    vi: {x: 270*Math.cos(3*Math.PI/4)+400, y: 330-270*Math.sin(3*Math.PI/4)},
    v: {x: 270*Math.cos(Math.PI)+400, y: 330-270*Math.sin(Math.PI)},
    iv: {x: 270*Math.cos(5*Math.PI/4)+400, y: 330-270*Math.sin(5*Math.PI/4)},
    iii: {x: 270*Math.cos(3*Math.PI/2)+400, y: 330-270*Math.sin(3*Math.PI/2)},
    ii: {x: 270*Math.cos(7*Math.PI/4)+400, y: 330-270*Math.sin(7*Math.PI/4)},
};

function divideLine(s, num, f =dotCords.dc, d=9) {
    num = +num;
    let l = num/(d-num);
    let x = (f.x + l*s.x)/(1+l);
    let y = (f.y + l*s.y)/(1+l);
    return {x: x, y: y}
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
        i: divideLine(dotCords.i, +data.i),
        ii: divideLine(dotCords.ii, +data.ii),
        iii: divideLine(dotCords.iii, +data.iii),
        iv: divideLine(dotCords.iv, +data.iv),
        v: divideLine(dotCords.v, +data.v),
        vi: divideLine(dotCords.vi, +data.vi),
        vii: divideLine(dotCords.vii, +data.vii),
        viii: divideLine(dotCords.viii, +data.viii),
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
            <td colspan="4" ><svg height="660" width="800" xmlns="http://www.w3.org/2000/svg">
                <g>
                    <rect height="660" width="800" y="0" x="0" stroke-width="1.5" stroke="#000" fill="#fff"/>
                    <circle cx="${dotCords.dc.x}" cy="${dotCords.dc.y}" r="270" stroke-width="1.5" stroke="#000" fill="#fff"/>
                    <circle cx="${dotCords.dc.x}" cy="${dotCords.dc.y}" r="210" stroke-width="1.5" stroke="#000" fill="#fff"/>
                    <circle cx="${dotCords.dc.x}" cy="${dotCords.dc.y}" r="120" stroke-width="1.5" stroke="#000" fill="#fff"/>
           
                     <path d="M ${dotCords.i.x} ${dotCords.i.y}  L ${dotCords.v.x} ${dotCords.v.y}" fill="transparent" stroke="black"/>
                     <path d="M ${dotCords.viii.x} ${dotCords.viii.y}  L ${dotCords.iv.x} ${dotCords.iv.y}" fill="transparent" stroke="black"/>
                     <path d="M ${dotCords.vii.x} ${dotCords.vii.y}  L ${dotCords.iii.x} ${dotCords.iii.y}" fill="transparent" stroke="black"/>
                     <path d="M ${dotCords.vi.x} ${dotCords.vi.y}  L ${dotCords.ii.x} ${dotCords.ii.y}" fill="transparent" stroke="black"/>
                     
                     <circle cx="${dotCords.i.x}" cy="${dotCords.i.y}" r="5" fill="#000"/>
                     <circle cx="${dotCords.viii.x}" cy="${dotCords.viii.y}" r="5" fill="#000"/>
                     <circle cx="${dotCords.vii.x}" cy="${dotCords.vii.y}" r="5" fill="#000"/>
                     <circle cx="${dotCords.vi.x}" cy="${dotCords.vi.y}" r="5" fill="#000"/>
                     <circle cx="${dotCords.v.x}" cy="${dotCords.v.y}" r="5" fill="#000"/>
                     <circle cx="${dotCords.iv.x}" cy="${dotCords.iv.y}" r="5" fill="#000"/>
                     <circle cx="${dotCords.iii.x}" cy="${dotCords.iii.y}" r="5" fill="#000"/>
                     <circle cx="${dotCords.ii.x}" cy="${dotCords.ii.y}" r="5" fill="#000"/>
                     
                     <circle cx="${dotCords.dc.x}" cy="${dotCords.dc.y}" r="5" fill="#000"/>
                     
                     <circle cx="${resDotCords.i.x}" cy="${resDotCords.i.y}" r="5" fill="#f00"/>
                     <circle cx="${resDotCords.ii.x}" cy="${resDotCords.ii.y}" r="5" fill="#f00"/>
                     <circle cx="${resDotCords.iii.x}" cy="${resDotCords.iii.y}" r="5" fill="#f00"/>
                     <circle cx="${resDotCords.iv.x}" cy="${resDotCords.iv.y}" r="5" fill="#f00"/>
                     <circle cx="${resDotCords.v.x}" cy="${resDotCords.v.y}" r="5" fill="#f00"/>
                     <circle cx="${resDotCords.vi.x}" cy="${resDotCords.vi.y}" r="5" fill="#f00"/>
                     <circle cx="${resDotCords.vii.x}" cy="${resDotCords.vii.y}" r="5" fill="#f00"/>
                     <circle cx="${resDotCords.viii.x}" cy="${resDotCords.viii.y}" r="5" fill="#f00"/>
                     
                     <polygon points="${resDotCords.i.x},${resDotCords.i.y} ${resDotCords.ii.x},${resDotCords.ii.y} ${resDotCords.iii.x},${resDotCords.iii.y}  ${resDotCords.iv.x},${resDotCords.iv.y} 
                     ${resDotCords.v.x},${resDotCords.v.y} ${resDotCords.vi.x},${resDotCords.vi.y} ${resDotCords.vii.x},${resDotCords.vii.y}  ${resDotCords.viii.x},${resDotCords.viii.y}"
                     style="fill:rgba(97,168,255,0.4);stroke:red;stroke-width:1;fill-rule:evenodd;" />
                     
                     
                     <text x="${dotCords.i.x+2}" y="${dotCords.i.y-10}"  font-size="14" font-family="arial, sans-serif">i. Экстраверсия</text>
                     <text x="${dotCords.ii.x}" y="${dotCords.ii.y-10}"  font-size="14" font-family="arial, sans-serif">ii. Спонтанность</text>
                     <text x="${dotCords.iii.x-60}" y="${dotCords.iii.y+17}"  font-size="14" font-family="arial, sans-serif">iii. Агрессивность</text>
                     <text x="${dotCords.iv.x-110}" y="${dotCords.iv.y-10}"  font-size="14" font-family="arial, sans-serif">iv. Ригидность</text>
                     <text x="${dotCords.v.x-100}" y="${dotCords.v.y-10}"  font-size="14" font-family="arial, sans-serif">v. Интроверсия</text>
                     <text x="${dotCords.vi.x-110}" y="${dotCords.vi.y-10}"  font-size="14" font-family="arial, sans-serif">vi. Сензитивность</text>
                     <text x="${dotCords.vii.x-60}" y="${dotCords.vii.y-10}"  font-size="14" font-family="arial, sans-serif">vii. Тревожность</text>
                     <text x="${dotCords.viii.x}" y="${dotCords.viii.y-10}"  font-size="14" font-family="arial, sans-serif">viii. Эмотивность</text>
                     
                     <text x="${dotCords.i.x}" y="${dotCords.i.y+15}"  font-size="14" font-style="italic" font-family="arial, sans-serif">(Соц. активность)</text>
                     <text x="${dotCords.dc.x-60}" y="${dotCords.dc.y+15}"  font-size="14" font-style="italic" font-family="arial, sans-serif">(Нормативность)</text>
                     <text x="${dotCords.dc.x+3}" y="${dotCords.dc.y-10}"  font-size="14" font-family="arial, sans-serif">0</text>
                     <text x="${dotCords.dc.x+3}" y="${divideLine(dotCords.vii, 4).y-6}"  font-size="14" font-family="arial, sans-serif">4</text>
                     <text x="${dotCords.dc.x+3}" y="${divideLine(dotCords.vii, 7).y-6}"  font-size="14" font-family="arial, sans-serif">7</text>
                     <text x="${dotCords.dc.x-45}" y="${dotCords.dc.y-24}"  font-size="14" font-family="arial, sans-serif">Стабильность</text>
                     <text x="${dotCords.v.x-120}" y="${dotCords.v.y+15}"  font-size="14" font-style="italic" font-family="arial, sans-serif">(Соц. пассивность)</text>
                     
                     <text x="${divideLine(dotCords.i, 1, dotCords.ii, 2).x+30}" y="${divideLine(dotCords.i, 1, dotCords.ii, 2).y}"  font-size="14" font-family="arial, sans-serif">Лидерство</text>
                     <text x="${divideLine(dotCords.ii, 1, dotCords.iii, 2).x+60}" y="${divideLine(dotCords.ii, 1, dotCords.iii, 2).y+15}"  font-size="14" font-family="arial, sans-serif">Некомформность</text>
                     <text x="${divideLine(dotCords.iii, 1, dotCords.iv, 2).x-160}" y="${divideLine(dotCords.iii, 1, dotCords.iv, 2).y+15}"  font-size="14" font-family="arial, sans-serif">Конфликтность</text>
                     <text x="${divideLine(dotCords.iv, 1, dotCords.v, 2).x-160}" y="${divideLine(dotCords.iv, 1, dotCords.v, 2).y}"  font-size="14" font-family="arial, sans-serif">Индивидуалистичность</text>
                     <text x="${divideLine(dotCords.v, 1, dotCords.vi, 2).x-110}" y="${divideLine(dotCords.v, 1, dotCords.vi, 2).y}"  font-size="14" font-family="arial, sans-serif">Зависимость</text>
                     <text x="${divideLine(dotCords.vi, 1, dotCords.vii, 2).x-135}" y="${divideLine(dotCords.vi, 1, dotCords.vii, 2).y-15}"  font-size="14" font-family="arial, sans-serif">Конформность</text>
                     <text x="${divideLine(dotCords.vii, 1, dotCords.viii, 2).x+40}" y="${divideLine(dotCords.vii, 1, dotCords.viii, 2).y-15}"  font-size="14" font-family="arial, sans-serif">Компромиссность</text>
                     <text x="${divideLine(dotCords.viii, 1, dotCords.i, 2).x+30}" y="${divideLine(dotCords.viii, 1, dotCords.i, 2).y}"  font-size="14" font-family="arial, sans-serif">Коммуникативность</text>


                     <text x="${resDotCords.i.x-3}" y="${resDotCords.i.y-6}"  font-size="14" fill="red" font-family="arial, sans-serif">${data.i}</text>
                     <text x="${resDotCords.ii.x-3}" y="${resDotCords.ii.y+18}"  font-size="14" fill="red" font-family="arial, sans-serif">${data.ii}</text>
                     <text x="${resDotCords.iii.x-15}" y="${resDotCords.iii.y+15}"  font-size="14" fill="red" font-family="arial, sans-serif">${data.iii}</text>
                     <text x="${resDotCords.iv.x-3}" y="${resDotCords.iv.y+18}"  font-size="14" fill="red" font-family="arial, sans-serif">${data.iv}</text>
                     <text x="${resDotCords.v.x-3}" y="${resDotCords.v.y-6}"  font-size="14" fill="red" font-family="arial, sans-serif">${data.v}</text>
                     <text x="${resDotCords.vi.x-3}" y="${resDotCords.vi.y-6}"  font-size="14" fill="red" font-family="arial, sans-serif">${data.vi}</text>
                     <text x="${resDotCords.vii.x-15}" y="${resDotCords.vii.y-6}"  font-size="14" fill="red" font-family="arial, sans-serif">${data.vii}</text>
                     <text x="${resDotCords.viii.x-3}" y="${resDotCords.viii.y-6}"  font-size="14" fill="red" font-family="arial, sans-serif">${data.viii}</text>                    

                </g>
            
            </svg></td>
        </tr>
        
        
        </table>`);

    container.append(newTbl)
}


