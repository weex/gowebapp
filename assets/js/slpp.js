var payment_hash = "";
var period = 5;
var expiration = current_time = 0;
var theTimer = null;

var HttpClient = function() {
    this.get = function(aUrl, aCallback) {
        var anHttpRequest = new XMLHttpRequest();
        anHttpRequest.onreadystatechange = function() {
            if (anHttpRequest.readyState == 4 && anHttpRequest.status == 200)
                aCallback(anHttpRequest.response);
        }

        anHttpRequest.open( "GET", aUrl, true );
        anHttpRequest.send( null );
    }
}
var client = new HttpClient();

function myTimer() {
    client.get('/check_invoice?payment_hash='+payment_hash, function(response) {
        var obj = JSON.parse(response)
        expiration = obj.expiry + obj.creation_date;
        if (current_time == 0)
            current_time = obj.creation_date + 10; // 10 secs for margin of safety
        current_time += period;
        if (current_time > expiration)
            timerStop("Invoice expired, please try again.");

        if (obj.settled === true)
            timerStop("Payment complete.");
    });
}

function timerStop(message) {
    clearInterval(theTimer);
    document.getElementById("invoice-qr").innerHTML = message;
    document.getElementById("invoice-text").innerHTML = '';
}

function getInvoice() {
    client.get('/invoice', function(response) {
        var obj = JSON.parse(response)
        payment_hash = obj.payment_hash;
        new QRCode(document.getElementById("invoice-qr"), "lightning:" + obj.payment_request);
        document.getElementById("invoice-text").innerHTML = obj.payment_request;
    });
    theTimer = setInterval(myTimer, period * 1000);
}
