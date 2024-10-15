const main = document.querySelector(".history")
const nav = document.querySelector("nav")
const recordsContainer = document.querySelector(".records-container")
const loadMoreBtn = document.querySelector(".load-more-btn")
main.style.paddingTop = nav.clientHeight + "px"


window.onload = async function() {
	const quizID = main.getAttribute("quiz_id")
	if (!quizID) {
		recordsContainer.insertAdjacentHTML("beforeend", "<h2>No history data found</h2>")
		return false
	}

	let lastQuizID = parseInt(quizID)
	console.log(lastQuizID)
	let quizzes = await fetchQuizzes(lastQuizID + 1)
	lastQuizID = quizzes[quizzes.length - 1].ID
	console.log(lastQuizID)
	appendRecords(recordsContainer, quizzes)

	loadMoreBtn.onclick = async function() {
		if (!quizzes || quizzes.length < 5) {
			this.innerHTML = "No more data"
			return false
		}

		quizzes = await fetchQuizzes(lastQuizID)
		console.log(quizzes)
		lastQuizID = quizzes[quizzes.length - 1].ID
		console.log(lastQuizID)
		appendRecords(recordsContainer, quizzes)
	}

}

function appendRecords(dest, quizzes) {
	for (let quiz of quizzes) {
		console.log(quiz)
		const recordNode = createRecordsNode(Object.entries(quiz))
		dest.append(recordNode)
	}
}


function createRecordsNode(obj) {
	const recordsNode = document.createElement("ul")
	recordsNode.className = "records"

	for (let [k, v] of obj) {
		switch (k) {
			case "CreatedAt":
				recordsNode.insertAdjacentHTML("beforeend", `<li><h3>Date: ${v}</h3></li>`)
				break
			case "Type":
				recordsNode.insertAdjacentHTML("beforeend", `<li><h3>Type: ${v}</h3></li>`)
				break
			case "Category":
				recordsNode.insertAdjacentHTML("beforeend", `<li><h3>Category: ${v}</h3></li>`)
				break
			case "ID":
				recordsNode.insertAdjacentHTML("beforeend", `<li class="trivias" quiz-id="${v}"><h3>Trivias: ▷</h3></li>`)
				break
			case "ScoreID":
				recordsNode.insertAdjacentHTML("beforeend", `<li class="score" score-id="${v}"><h3>Score: ▷</h3></li>`)
				break
		}
	}

	let itemTrivias = recordsNode.querySelector(".trivias")
	const quizID = itemTrivias.getAttribute("quiz-id")
	console.log(quizID)
	if (quizID != null) {
		itemTrivias.onclick = async function(event) {
			if (event.target != this.firstElementChild || event.defaultPrevented) {
				return false
			}
			if (this.childNodes.length > 1) {
				if (!this.lastElementChild.hidden) {
					this.lastElementChild.hidden = true
				} else {
					this.lastElementChild.hidden = false
				}
				return false
			}
			let trivias = await fetchTrivias(quizID)
			console.log(trivias)
			appendTriviasContainer(this, trivias, quizID)
		}
	}

	let itemScore = recordsNode.querySelector(".score")
	const scoreID = itemScore.getAttribute("score-id")
	console.log(scoreID)
	if (scoreID != null) {
		itemScore.onclick = async function(event) {
			if (event.target != this.firstElementChild) {
				return false
			}
			if (this.childNodes.length > 1) {
				if (!this.lastElementChild.hidden) {
					this.lastElementChild.hidden = true
				} else {
					this.lastElementChild.hidden = false
				}
				return false
			}

			let response = await fetchScore(scoreID)
			if (response.Error) {
				this.innerHTML += " The quiz was not completed"
				this.onclick = ""
				return false
			}
			appendScore(this, Object.entries(response))
		}
	}


	return recordsNode
}

async function appendTriviasContainer(dest, arr, quizID) {
	const triviasContainer = document.createElement("div")
	triviasContainer.className = "trivias-container"
	const answers = await fetchAnswers(quizID)
	if (answers == null) {
		dest.innerHTML += " The quiz was not completed"
		dest.onclick = ""
		return
	}
	let triviasID = []
	let answersID = []
	for (let answer of answers) {
		for (let [k, v] of Object.entries(answer)) {
			if (!answer) {
				continue
			}
			switch (k) {
				case "ID":
					answersID.push(v)
					break
				case "TriviaID":
					triviasID.push(v)
					break
			}
		}
	}
	console.log(answers)
	let i = 1
	for (let obj of arr) {
		const trivia = document.createElement("ul")
		trivia.className = "trivia"
		triviasContainer.insertAdjacentHTML("beforeend", `<h3>#${i}</h3>`)
		i += 1
		for (let [k, v] of Object.entries(obj)) {
			switch (k) {
				case "Category":
					trivia.insertAdjacentHTML("beforeend", `<li><h4>Category: ${v}</h4></li>`)
					break
				case "Difficulty":
					trivia.insertAdjacentHTML("beforeend", `<li><h4>Category: ${v}</h4></li>`)
					break
				case "Type":
					trivia.insertAdjacentHTML("beforeend", `<li><h4>Type: ${v}</h4></li>`)
					break
				case "ID":
					const li = document.createElement("li")
					appendOption(li, answersID[triviasID.indexOf(v)])
					trivia.appendChild(li)
					break
				case "Question":
					trivia.insertAdjacentHTML("beforeend", `<li><h4>Question: ${v}</h4></li>`)
					break
			}
		}

		triviasContainer.appendChild(trivia)
	}

	dest.append(triviasContainer)
}


async function appendOption(dest, answerID) {
	const option = await fetchOption(answerID)
	const optionNode = document.createElement("h4")
	for (let [k, v] of Object.entries(option)) {
		switch (k) {
			case "Option":
				optionNode.innerHTML = `Answer: ${v}`
				break
			case "Correct":
				if (v) {
					optionNode.innerHTML += " ✔️"
				} else {
					optionNode.innerHTML += " ✖️"
				}
				break
		}
	}
	dest.append(optionNode)
}


function appendScore(dest, obj) {
	const scoreNode = document.createElement("ul")
	scoreNode.className = "score-list"
	for (let [k, v] of obj) {
		switch (k) {
			case "CompletionTime":
				scoreNode.insertAdjacentHTML("beforeend", `<li><h3>Completion time: ${v}</h3></li>`)
				break
			case "TotalScore":
				scoreNode.insertAdjacentHTML("beforeend", `<li><h3>Score: ${v}</h3></li>`)
				break
		}
	}

	dest.append(scoreNode)
}


async function fetchOption(answerID) {
	const response = await fetch(`/user/history/option?id=${answerID}`,
		{
			method: "GET",
			headers: {
				"Content-type": "application/json"
			}
		}
	)
	const option = await response.json()
	return option
}

async function fetchAnswers(quizID) {
	const response = await fetch(`/user/history/answer?id=${quizID}`,
		{
			method: "GET",
			headers: {
				"Content-type": "application/json"
			}
		}
	)
	if (!response.ok) {
		return {
			Error: "The quiz was not answered"
		}
	}
	const answers = await response.json()
	return answers
}

async function fetchTrivias(quizID) {
	const response = await fetch(`/user/history/trivia?id=${quizID}`,
		{
			method: "GET",
			headers: {
				"Content-type": "application/json"
			}
		}
	)
	if (!response.ok) {
		return {
			Error: "The quiz was not completed"
		}
	}
	const trivias = await response.json()
	return trivias
}

async function fetchScore(scoreID) {
	const response = await fetch(`/user/history/score?id=${scoreID}`,
		{
			method: "GET",
			headers: {
				"Content-type": "application/json"
			}
		}
	)
	if (!response.ok) {
		return {
			Error: "The quiz was not completed"
		}
	}
	const score = await response.json()
	return score
}

async function fetchQuizzes(quizID) {
	const response = await fetch(`/user/history/load?id=${quizID}`,
		{
			method: "GET",
			headers: {
				"Content-type": "application/json"
			}
		}
	)
	if (!response.ok) {
		return {
			Error: "No history data found"
		}
	}
	const quizzes = await response.json()
	return quizzes
}
