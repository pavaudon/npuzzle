var socket = null;
var GlobalMap = null;
var nbMoves = 0;
var Conf = 0;
var Allmaps = [];
var sleeptime = 250;
var moves, index, myInterval;
$(window).on('load', function () {
    socket = new WebSocket("ws://127.0.0.1:8080/event_ws");
    socket.onopen = function (e) {
        console.log("CONNECTED")
        let evt = {
            command: "get_conf",
            value: ""
        }
        socket.send(JSON.stringify(evt))
    };

    socket.onmessage = function (event) {
        let evt = JSON.parse(event.data)
        console.log(evt)
        switch (evt.command) {
            case "map":
                Allmaps.unshift(evt.value);
                createMap(JSON.parse(evt.value).map);
                if (Conf.StartAfterGen) {
                    start()
                }
                break
            case "set_conf":
                changeTheConf(JSON.parse(evt.value))
                break
            case "solved":
                $(".btn-secondary, .btn-info, .btn-light").prop('disabled', true);
                $("#start").addClass("btn-warning").html("PAUSE")
                $("#loading").addClass("d-none");
                $("#cost").html(JSON.parse(evt.value).Cost)
                $("#heuristic").html(JSON.parse(evt.value).Heuristic)
                $("#complexity-time").html(JSON.parse(evt.value).ComplexityTime)
                $("#complexity-size").html(JSON.parse(evt.value).ComplexitySize)
                moves = JSON.parse(evt.value).Moves;
                if (parseInt(JSON.parse(evt.value).Cost) === 0) {
                    alert("0 Move to resolt, ALREADY SOLVED")
                }
                index = 0;
                play()
                break
            case "timer":
                $("#timer").html(evt.value)
                break
            case "error":
                alert("GET ERROR FROM SERVER : " + evt.value)
                setLog(`[error] ` + evt.value);
                $("#loading").addClass("d-none");
                break
            case "ok":
                console.log("get OK")
                break
            default:
                setLog("[websocket] UNKNONW COMMAND")
        }
    };

    socket.onclose = function (evt) {
        console.log(evt)
        setLog(`[error] "WEBSOCKET CONNEXION IS CLOSED"`);
    };

    socket.onerror = function (evt) {
        console.log(evt)
        setLog(`[error] ${evt}`);
    };
});


function changeTheConf(conf) {
    console.log(conf)
    Conf = conf
    $("#sort-mode").val(conf.SortMode)
    $("#solvable").prop("checked", conf.Solvable)
    $("#map-size").val(conf.GenSize)
    $("#sag").prop("checked", conf.StartAfterGen)
    $("#map-iteration").val(conf.Iteration)
    $("#heuristic-mode").val(conf.Heuristic)
    $("#timer-chkb").prop("checked", conf.AlgoTimer)
}

function updateTheConf() {
    Conf.SortMode = parseInt($("#sort-mode").val())
    Conf.Solvable = $("#solvable").is(':checked')
    Conf.GenSize = parseInt($("#map-size").val())
    Conf.StartAfterGen = $("#sag").is(':checked')
    Conf.Iteration = parseInt($("#map-iteration").val())
    Conf.Heuristic = parseInt($("#heuristic-mode").val())
    Conf.AlgoTimer = $("#timer-chkb").is(':checked')
    let evt = {
        command: "set_conf",
        value: JSON.stringify(Conf)
    }
    socket.send(JSON.stringify(evt))
}


function importMap(snailed) {

    if (snailed) {
        $("#sort-mode").val(1)
    } else {
        $("#sort-mode").val(0)
    }
    updateTheConf()
    let evt = {
        command: "import_map",
        value: $("#import-value").val()
    }
    socket.send(JSON.stringify(evt))
    $('#import-modal').modal('hide');
}

function generateMap() {
    updateTheConf()
    let evt = {
        command: "generate_map",
        value: ""
    }
    socket.send(JSON.stringify(evt))
}

function start() {
    if ($(".btn-secondary").is(":disabled")) {
        if ($("#start").html() === "START") {
            $("#start").addClass("btn-warning").html("PAUSE")
            play()
        } else {
            $("#start").removeClass("btn-warning").html("START")
            clearInterval(myInterval)
        }
        return
    }
    updateTheConf()
    if (nbMoves > 0) {
        if (confirm("do you want to start the Algo with the current map or with the initial map ? (Ok to choose the current and Cancel to choose the initial)")) {
            let evt = {
                command: "import_map",
                value: getMapFromHTML()
            }
            socket.send(JSON.stringify(evt));
        }
    }
    nbMoves = 0
    $("#nb-moves").html(0)
    $("#loading").removeClass("d-none");
    let evt = {
        command: "start",
        value: ""
    }
    socket.send(JSON.stringify(evt));
}

function ImportMapFromHistoric(index) {

        createMap(JSON.parse(Allmaps[index]).map)
        let evt = {
            command: "import_map",
            value: getMapFromHTML()
        }
        socket.send(JSON.stringify(evt));
        $('#maps-modal').modal('hide');

}

function Historic() {
    $("#all-maps").html("")
    for (let j = 0; j < Allmaps.length; j++) {
        let CurrentMap = JSON.parse(Allmaps[j]).map
        let size = (100 / CurrentMap.length)
        let allCases = ``
        for (let i = 0; i < CurrentMap.length; i++) {
            for (let ii = 0; ii < CurrentMap.length; ii++) {
                if (CurrentMap[i][ii] === "0") {
                    allCases += `<button type="button" class="border taq-0 btn btn-light btn-lg p-1" style="width: ` + size + `%">0</button>`
                } else {
                    allCases += `<button type="button" class="taq-` + CurrentMap[i][ii] + ` border btn btn-info btn-lg p-1" style="width: ` + size + `%">` + CurrentMap[i][ii] + `</button>`
                }
            }
        }
        $("#all-maps").append(`<div class=" border rounded m-1 p-1 text-center" style="width: `+ parseInt(CurrentMap.length * 64 + 18)+`px"><div><button class="btn btn-secondary my-1" onclick="ImportMapFromHistoric(`+j+`)">Choose this map</button></div><div>`+allCases+`</div></div>`)
    }
    $("#maps-modal").modal("show")
}

function createMap(map) {
    GlobalMap = map
    $("#my-map").css("width", ((map.length * 64) + 18) + "px").html("")
    nbMoves = 0
    $("#nb-moves").html(nbMoves)
    let size = (100 / map.length)
    $("#map-size").val(map.length)
    for (let i = 0; i < map.length; i++) {
        for (let ii = 0; ii < map.length; ii++) {
            if (map[i][ii] === "0") {
                $("#my-map").append('<button type="button" class="border taq-0 btn btn-light btn-lg p-3" style="width: ' + size + '%"></button>')
            } else {
                $("#my-map").append('<button type="button" class="taq-' + map[i][ii] + ' border btn btn-info btn-lg p-3" style="width: ' + size + '%" onclick="makeAMove(\'' + map[i][ii] + '\')" >' + map[i][ii] + '</button>')
            }
        }
    }
}

function getMapFromHTML() {
    let size = Math.sqrt($("#my-map button").length)
    let map = size + "\n"
    let i = 1
    $("#my-map button").each(function () {
        if (this.innerHTML.length == 0) {
            map += "0"
        } else {
            map += this.innerHTML
        }
        if (i % size === 0) {
            map += "\n"
        } else {
            map += " "
        }
        i++;
    })
    return map
}


function makeAMove(value) {
    for (let i = 0; i < GlobalMap.length; i++) {
        for (let ii = 0; ii < GlobalMap.length; ii++) {
            if (GlobalMap[i][ii] === value) {
                checkTheMove(i, ii)
                return
            }
        }
    }
}

function checkTheMove(x, y) {
    if (x > 0 && GlobalMap[x - 1][y] === "0") { //Au dessus
        DoTheMove(x - 1, y, x, y, GlobalMap[x][y])
    } else if (x < GlobalMap.length - 1 && GlobalMap[x + 1][y] === "0") { //Au dessous
        DoTheMove(x + 1, y, x, y, GlobalMap[x][y])
    } else if (y > 0 && GlobalMap[x][y - 1] === "0") { // A GAUCHE
        DoTheMove(x, y - 1, x, y, GlobalMap[x][y])
    } else if (y < GlobalMap.length - 1 && GlobalMap[x][y + 1] === "0") { //A DROITE
        DoTheMove(x, y + 1, x, y, GlobalMap[x][y])
    } else {
        setLog("UNABLE MOVE")
    }
}

function DoTheMove(zeroX, zeroY, x, y, value) {
    nbMoves++;
    $("#nb-moves").html(nbMoves)
    GlobalMap[zeroX][zeroY] = value;
    GlobalMap[x][y] = "0";
    $(".taq-" + value).addClass("taq-x btn-light").removeClass("btn-info taq-" + value).html('').attr("onclick", "");
    $(".taq-0").addClass("taq-" + value + " btn-info").removeClass("btn-light taq-0").html(value).attr("onclick", "makeAMove(\'" + value + "\')");
    $(".taq-x").addClass("taq-0").removeClass("taq-x")
}

function play() {

    sleeptime = parseInt($("#sleep-time").val())
    myInterval = setInterval(function () {
        if (!moves || index >= moves.length) {
            clearInterval(myInterval)
            $(".btn-secondary,.btn-info, .btn-light ").prop('disabled', false);
            $("#start").removeClass("btn-warning").html("START")
            return
        }

        let coods = moves[index].split("x")
        checkTheMove(parseInt(coods[0]), parseInt(coods[1]))
        index++;

    }, sleeptime);
}



function setLog(text) {
    $("#log").html($("#log").html() + text + " ; ")
}