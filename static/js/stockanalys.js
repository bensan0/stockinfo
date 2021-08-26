

function getCorpDaily() {
    let date = document.getElementById('cpdailyyyyymmdd').value;
    console.log(date)
    let req = new Request('http://127.0.0.1:8080/stockinfo/crawlers/corpdaily?date=' + date);
    fetch(req)
        .then(res => {
            if (res.status == 200) {
                return res.json()
            } else {
                throw new Error(res.statusText)
            }
        })
        .then(jsonObj => {
            console.log(jsonObj)
            let err = jsonObj['error']
            if (err === null || err.length === 0) {
                document.getElementById('cpdailystatus').innerHTML = 'OK';
            } else {
                document.getElementById('cpdailystatus').innerHTML = err;
            }
        })
        .catch(error => {
            document.getElementById('cpdailystatus').innerHTML = error;
        });
}

function getStockDaily() {
    let date = document.getElementById('stdailyyyyymmdd').value;
    console.log(date)
    let req = new Request('http://127.0.0.1:8080/stockinfo/crawlers/dailyquot?date=' + date);
    fetch(req)
        .then(res => {
            if (res.status == 200) {
                return res.json()
            } else {
                throw new Error(res.statusText)
            }
        })
        .then(jsonObj => {
            console.log(jsonObj)
            let err = jsonObj['error']
            if (err === null || err.length === 0) {
                document.getElementById('stdailystatus').innerHTML = 'OK';
            } else {
                document.getElementById('stdailystatus').innerHTML = err;
            }
        })
        .catch(error => {
            document.getElementById('stdailystatus').innerHTML = error;
        });
}

function getDistribution() {
    let req = new Request('http://127.0.0.1:8080/stockinfo/crawlers/distribution');
    fetch(req)
        .then(res => {
            if (res.status == 200) {
                return res.json()
            } else {
                throw new Error(res.statusText)
            }
        })
        .then(jsonObj => {
            console.log(jsonObj)
            let err = jsonObj['error']
            if (err === null || err.length === 0) {
                document.getElementById('distributionstatus').innerHTML = 'OK';
            } else {
                document.getElementById('distributionstatus').innerHTML = err;
            }
        })
        .catch(error => {
            document.getElementById('distributionstatus').innerHTML = error;
        });
}

function filtTradingvol() {
    clearFilter();
    let date = document.getElementById('filterdate').value
    let percent = document.getElementById('filterpercent').value
    let lowervol = document.getElementById('filterlowerlimit').value
    let highervol = document.getElementById('filterhigherlimit').value
    let req = new Request('http://127.0.0.1:8080/stockinfo/stockdaily/filttradingvol?' + 'date=' + date + '&percent=' + percent + '&lowervol=' + lowervol + '&highervol=' + highervol /*+ '&fluc=' + fluc*/);

    fetch(req)
        .then(res => {
            if (res.status == 200) {
                return res.json()
            } else {
                throw new Error(res.statusText)
            }
        })
        .then(jsonObj => {
            console.log(jsonObj);
            let err = jsonObj['error']
            if (err === null || err.length === 0) {
                document.getElementById('filterstatus').innerHTML = 'OK';
                let filter = document.getElementById('filtertbody');
                let filtertranstbody = document.getElementById('filtertranstbody');
                jsonObj['data'].forEach(function (data) {
                    filter.innerHTML += '<tr>' +
                        '<th scope="row">' + data['code'] + '</th>' +
                        '<td>' + data['name'] + '</td>' +
                        '<td>' + data['industry'] + '</td>' +
                        '</tr>';

                    data['trans'].forEach(function (datatrans) {
                        filtertranstbody.innerHTML += '<tr>' +
                            '<th scope="row">' + datatrans['code'] + '</th>' +
                            '<td>' + datatrans['name'] + '</td>' +
                            '<td>' + datatrans['tradingVol'] + '</td>' +
                            '<td>' + datatrans['closing'] + '</td>' +
                            '<td>' + datatrans['flucPercent'] + '</td>' +
                            '<td>' + datatrans['fluctuation'] + '</td>' +
                            '<td>' + datatrans['date'] + '</td>' +
                            '</tr>';
                    });
                });
            } else {
                document.getElementById('filterstatus').innerHTML = err;
            }
        })
        .catch(error => {
            document.getElementById('distributionstatus').innerHTML = error;
        });
}

function getStockTrans() {
    clearStockTrans();
    let days = document.getElementById('stocktransdays').value
    let code = document.getElementById('stocktranscode').value
    let req = new Request('http://127.0.0.1:8080/stockinfo/stockdaily/getdays?' + 'days=' + days + '&code=' + code);

    fetch(req)
        .then(res => {
            if (res.status == 200) {
                return res.json()
            } else {
                throw new Error(res.statusText)
            }
        })
        .then(jsonObj => {
            console.log(jsonObj);
            let err = jsonObj['error']
            if (err === null || err.length === 0) {
                document.getElementById('stocktransstatus').innerHTML = 'OK';
                let stocktranstbody = document.getElementById('stocktranstbody');
                jsonObj['data'].forEach(function (data) {
                    stocktranstbody.innerHTML += '<tr>' +
                        '<th scope="row">' + data['code'] + '</th>' +
                        '<td>' + data['name'] + '</td>' +
                        '<td>' + data['tradingVol'] + '</td>' +
                        '<td>' + data['closing'] + '</td>' +
                        '<td>' + data['flucPercent'] + '</td>' +
                        '<td>' + data['fluctuation'] + '</td>' +
                        '<td>' + data['date'] + '</td>' +
                        '</tr>';
                });
            } else {
                document.getElementById('filterstatus').innerHTML = err;
            }
        })
        .catch(error => {
            document.getElementById('distributionstatus').innerHTML = error;
        });
}

function getCorpTrans() {
    clearCorpTrans();
    let days = document.getElementById('corptransdays').value
    let code = document.getElementById('corptranscode').value

    let req = new Request('http://127.0.0.1:8080/stockinfo/corpdaily/getdays?code=' + code + '&days=' + days);

    fetch(req)
        .then(res => {
            if (res.status == 200) {
                return res.json()
            } else {
                throw new Error(res.statusText)
            }
        })
        .then(jsonObj => {
            console.log(jsonObj);
            let fiTotal = 0, fcTotal = 0, iTTotal = 0, dTotal = 0, dSTotal = 0, dHTotal = 0, tTotal = 0;


            let err = jsonObj['error']
            if (err === null || err.length === 0) {

                document.getElementById('corptransstatus').innerHTML = 'OK';
                const corptranstbody = document.getElementById('corptranstbody');
                jsonObj['data'].forEach(function (data) {
                    corptranstbody.innerHTML += '<tr>' +
                        '<th scope="row">' + data['code'] + '</th>' +
                        '<td>' + data['name'] + '</td>' +
                        '<td>' + data['foreignInvestors'] + '</td>' +
                        '<td>' + data['foreignCorp'] + '</td>' +
                        '<td>' + data['investmentTrust'] + '</td>' +
                        '<td>' + data['dealer'] + '</td>' +
                        '<td>' + data['dealerSelf'] + '</td>' +
                        '<td>' + data['dealerHedge'] + '</td>' +
                        '<td>' + data['total'] + '</td>' +
                        '<td>' + data['date'] + '</td>' +
                        '</tr>';
                    fiTotal += data['foreignInvestors'];
                    fcTotal += data['foreignCorp'];
                    iTTotal += data['investmentTrust'];
                    dTotal += data['dealer'];
                    dSTotal += data['dealerSelf'];
                    dHTotal += data['dealerHedge'];
                    tTotal += data['total'];
                });
                corptranstbody.innerHTML += '<tr>' +
                    '<th scope="row" colspan="2">合計</th>' +
                    '<td>' + fiTotal + '</td>' +
                    '<td>' + fcTotal + '</td>' +
                    '<td>' + iTTotal + '</td>' +
                    '<td>' + dTotal + '</td>' +
                    '<td>' + dSTotal + '</td>' +
                    '<td>' + dHTotal + '</td>' +
                    '<td>' + tTotal + '</td>' +
                    '<td></td>' +
                    '</tr>';
            } else {
                document.getElementById('filterstatus').innerHTML = err;
            }
        })
        .catch(error => {
            document.getElementById('corptransstatus').innerHTML = error
        })
}

function clearFilter() {
    document.getElementById('filtertbody').innerHTML = '';
    document.getElementById('filtertranstbody').innerHTML = '';
    document.getElementById('filterstatus').innerHTML = '';
}

function clearStockTrans() {
    document.getElementById('stocktranstbody').innerHTML = '';
    document.getElementById('stocktransstatus').innerHTML = '';
}

function clearCorpTrans() {
    document.getElementById('corptranstbody').innerHTML = '';
    document.getElementById('corptransstatus').innerHTML = '';
}