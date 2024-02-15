const ps: HTMLCollectionOf<HTMLParagraphElement> = document.getElementsByTagName("p");
for (let i = 0; i < ps.length; i++) {
    ps[i].style.color = "red";
    ps[i].innerText = "Hello, World!";
}