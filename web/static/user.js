$(document).ready(function() {

	$('.btn-logout').click(function(e) {
		Cookies.remove('auth-session');
	});
});
document.querySelector("main").style.height = document.documentElement.clientHeight + "px"

document.querySelector("form").clientHeight
