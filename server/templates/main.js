window.addEventListener('DOMContentLoaded', function() {
	console.log("DOM loaded");
	document.getElementById("login-btn").addEventListener("click", auth);

	document.getElementById("logout-form").style.display = "none";

	is_auth = checkAuth();

	if (is_auth) {
		console.log("is_auth", is_auth);
		document.getElementById("login-form").style.display = "none";
		document.getElementById("logout-form").style.display = "block";
	}

	if (!is_auth && window.location.pathname != "/login") {
		console.log("redirect to login");
		window.location.href = "/login";
	}

	else {
		console.log("no redirect");
	}
});

async function auth() {
	response = await fetch("/oauth/token", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
			// "HTMX-Request": "true",
		},
		body: JSON.stringify({
			username: document.querySelector("input[name='username']").value,
			password: document.querySelector("input[name='password']").value,
			grant_type: "password",
		}),
	});

	// console.log("response", response);
	data = await response.json();

	console.log("data", data);
	document.cookie = `jwt=${data.access_token}; path=/;`;
	document.cookie = `refresh_token=${data.refresh_token}; path=/;`;
	location.reload();
}

function checkAuth() {
	// console.log("checkAuth");
	var jwt = getCookie("jwt");
	if (jwt) {
		console.log("jwt", jwt);
		return true;
	} else {
		console.log("no jwt");
		return false;
	}
}

function getCookie(name) {
	var value = "; " + document.cookie;
	var parts = value.split("; " + name + "=");
	if (parts.length == 2)
		return parts
			.pop()
			.split(";")
			.shift();
}

function deleteCookie(name) {
	document.cookie = name + "=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;";
}
