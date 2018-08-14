var data = [
    {{range.Payments}}
        {
            type: 'Payment',
            date: '{{.Date.Format "Jan 2, 2006"}}',
            balchg: '{{.Amount}}'
        },
    {{end}}
    {{range .Cuts}}
        {
            type: 'Cut',
            date: '{{.Date.Format "Jan 2, 2006"}}',
            balchg: '-{{.Price}}'
        },
    {{end}}
];

var sorted = data.sort(function (a, b) {
    return new Date(b.date) - new Date(a.date);
});

var timeline = document.getElementById("timeline");
for (i = sorted.length - 1; i >= 0; i--) {
    var r = timeline.insertRow(1);
    r.insertCell(0).innerHTML = sorted[i].type;
    r.insertCell(1).innerHTML = sorted[i].date;
    r.insertCell(2).innerHTML = "$" + sorted[i].balchg;

    var color = "rgba(255, 255, 255, 0)";
    if (sorted[i].balchg >= 0) {
        color = "rgba(50, 255, 100, .75)";
    } else {
        color = "rgba(255, 100, 75, .75)";
    }
    r.style.backgroundColor = color;
}