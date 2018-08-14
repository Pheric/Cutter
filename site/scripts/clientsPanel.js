let selectedClient = "";
let page = "summary";

function selectClient(c) {
    if(selectedClient === c) return;

    if(selectedClient !== '') document.getElementById(selectedClient).classList.remove("selected");
    document.getElementById(c).classList.add("selected");
    selectedClient = c;

    changePage(page, true);
}

function changePage(pg, force) {
    if(pg === page && !force) return;

    let xhttp = new XMLHttpRequest();
    let callback = function(t){}; // Can't use nil :(

    switch(pg) {
        case "summary":
            xhttp.open("GET", "getCbSummary?cid=" + selectedClient, true);
            callback = function(t) {
                document.getElementById("cb").innerHTML = t.responseText;

                loadSummary();
            };

            document.getElementById("cb-nav-" + page).style.borderBottomColor = "rgba(0, 0, 0, 0)";
            page = pg;
            document.getElementById("cb-nav-" + pg).style.borderBottomColor = "green";
            break;
        case "log":
            xhttp.open("GET", "getCbLogHtml", true);
            callback = function(t) {
                document.getElementById("cb").innerHTML = t.responseText;
                changePage("log_", false);
            };

            document.getElementById("cb-nav-" + page).style.borderBottomColor = "rgba(0, 0, 0, 0)";
            page = pg;
            document.getElementById("cb-nav-" + pg).style.borderBottomColor = "green";
            break;
        case "log_":
            xhttp.open("GET", "getCbLogJs?cid=" + selectedClient, true);
            callback = function(t) {
                eval(t.responseText);
            };
            break;
        case "edit":
            xhttp.open("GET", "cbEdit?cid=" + selectedClient, true);
            callback = function(t) {
                document.getElementById("cb").innerHTML = t.responseText;
                loadEditFields()
            };
            break;
    }

    xhttp.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200){
            callback(this)
        }
    };
    xhttp.send();
}

function loadSummary() {
    let pd = "";
    let p = document.getElementById("period").innerText;
    let q = document.getElementById("quote").innerText;
    let t = document.getElementById("ttc").innerText;

    switch(p){
    case 0:
        pd = "Weekly";
        break;
    case 1:
        pd = "Bi-Weekly";
        break;
    case 2:
        pd = "On Demand";
        break;
    }
    document.getElementById("pd").innerText = pd;

    document.getElementById("profitData").innerText = "(profit of ~ $" + (q / t) + " / minute)";
}

function loadEditFields() {
    let period = document.getElementById("period").innerText;

    document.getElementById("pd" + period).checked = true;
    if(document.getElementById("uuid").innerText !== ""){
        document.getElementById("jsubmit").disabled = false;
    }
}