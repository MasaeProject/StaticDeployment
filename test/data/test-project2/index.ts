window.addEventListener("load", () => {
  const ps: HTMLCollectionOf<HTMLParagraphElement> =
    document.getElementsByTagName("p");
  for (let i = 0; i < ps.length; i++) {
    const p: HTMLParagraphElement = ps[i];
    p.style.color = "orange";
    p.innerText = "Welcome!";
  }
});
