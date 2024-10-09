const main = document.querySelector("main")
main.style.top = document.querySelector("nav").clientHeight + "px"
document.body.style.height = document.documentElement.clientHeight - document.querySelector("form").clientHeight + 'px';
