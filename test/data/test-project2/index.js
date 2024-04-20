window.addEventListener("load", function () {
    var ps = document.getElementsByTagName("p");
    for (var i = 0; i < ps.length; i++) {
        var p = ps[i];
        p.style.color = "red";
        p.innerText = "ようこそ!";
    }
});
