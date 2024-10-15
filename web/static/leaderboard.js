const main = document.querySelector("main")
const nav = document.querySelector("nav")
const table = document.querySelector("table")
const loadMoreBtn = document.querySelector(".load-more-btn")
main.style.paddingTop = nav.clientHeight + "px"

let highestScore = parseInt(main.getAttribute("score"))
console.log(highestScore)

window.onload = async function() {
	let rows = await fetchScores(highestScore + 1)
	rows = Array.from(rows)
	highestScore = rows[rows.length - 1]
	console.log(highestScore)
	appendRows(table, rows)

	loadMoreBtn.onclick = async function() {
		if (!rows || rows.length < 5) {
			this.innerHTML = "No more data"
			return false
		}

		rows = await fetchScores(highestScore)
		console.log(rows)
		highestScore = rows[rows.length - 1]
		console.log(highestScore)
		appendRows(table, rows)
	}

}

function appendRows(dest, rows) {
	for (let row of rows) {
		const tr = document.createElement("tr")
		for (let [k, v] of Object.entries(row)) {
			switch (k) {
				case "Email":
					tr.insertAdjacentHTML("beforeend", `<td>${v}</td>`)
					break
				case "EasyQuizzesDone":
					tr.insertAdjacentHTML("beforeend", `<td>${v}</td>`)
					break
				case "MediumQuizzesDone":
					tr.insertAdjacentHTML("beforeend", `<td>${v}</td>`)
					break
				case "HardQuizzesDone":
					tr.insertAdjacentHTML("beforeend", `<td>${v}</td>`)
					break
				case "CompletionTime":
					tr.insertAdjacentHTML("beforeend", `<td>${v / 1000}</td>`)
					break
				case "Score":
					tr.insertAdjacentHTML("beforeend", `<td>${v}</td>`)
					break
			}
		}
		dest.append(tr)
	}
}

async function fetchScores(score) {
	const response = await fetch(`/leaderboard/load?score=${score}`,
		{
			method: "GET",
			headers: {
				"Content-type": "application/json"
			}
		}
	)
	const scores = await response.json()
	return scores
}
