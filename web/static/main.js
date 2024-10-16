const nav = document.querySelector("nav")
document.body.style.height = document.documentElement.clientHeight + nav.clientHeight + "px"
document.addEventListener("DOMContentLoaded", (event) => {
	document.body.addEventListener("htmx:beforeSwap", function(event) {
		if (event.detail.xhr.status === 422) {
			event.detail.shouldSwap = true;
			event.detail.isError = false;
		}
	})
})
