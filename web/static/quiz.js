document.querySelector("main").style.top = document.querySelector("nav").clientHeight + "px"
document.querySelector(".results-container").style.top = document.querySelector("nav").clientHeight + "px"

document.querySelector(".done-btn").addEventListener("click", function() {
	location.reload()
})
class Quizgo {
	constructor() {
		this.quizStore = document.querySelector(".quizzes")
		this.quizzes = document.querySelectorAll(".quiz")
		this.index = 0
		this.quizContainer = document.querySelector(".quiz-container")
		this.nextBtn = document.querySelector(".next-btn")
		this.results = document.querySelector(".results-container")
		this.quizID = this.quizContainer.getAttribute("quiz_id")
		this.startedAt = Date.now()
		this.quizzesDone = new Map([
			["easy", 0],
			["medium", 0],
			["hard", 0]
		]);

		this.quiz = this.quizzes[this.index]
		this.setUpOptions()

		this.quizContainer.append(this.quiz)
		this.nextBtn.onclick = () => {
			this.nextQuiz()
		}
	}
	shuffleOptions() {
		const optionsArray = Array.from(this.quiz.querySelectorAll(".options-list-item"))
		for (let i = optionsArray.length - 1; i > 0; i--) {
			const j = Math.floor(Math.random() * (i + 1));  // Random index from 0 to i
			[optionsArray[i], optionsArray[j]] = [optionsArray[j], optionsArray[i]];  // Swap
		}

		this.optionsList = this.quiz.querySelector(".options-list")
		this.optionsList.innerHTML = ""
		optionsArray.forEach(item => {
			this.optionsList.appendChild(item)
		})
	}

	setUpOptions() {
		this.shuffleOptions()
		this.options = this.optionsList.querySelectorAll(".option")
		this.options.forEach(option => {
			option.onchange = () => {
				this.quiz.setAttribute("selected", option.getAttribute("value"))
				this.quiz.setAttribute("selected-id", option.getAttribute("option-id"))
			}
		})
	}

	async storeResults() {
		let correctOptionValue = this.quiz.querySelector(".correct").getAttribute("value")
		let selectedOptionValue = this.quiz.getAttribute("selected")
		let difficulty = this.quiz.getAttribute("difficulty")

		let stats = document.createElement("h3")
		stats.innerText = `
Quiz #${this.index + 1}: ${this.quiz.querySelector(".question").innerHTML}
Type: ${this.quiz.getAttribute("type")}
Category: ${this.quiz.getAttribute("category")}
Difficulty: ${this.quiz.getAttribute("difficulty")}`

		let resultQuiz = document.createElement("ul")
		resultQuiz.className = "result-quiz"
		resultQuiz.append(stats)

		let optionValue = ""
		this.options.forEach(option => {
			optionValue = option.getAttribute("value")

			resultQuiz.insertAdjacentHTML("beforeend", `<li class="result-option">${optionValue}</li>`)

			this.results.append(resultQuiz)

			if (optionValue == correctOptionValue) {
				if (optionValue == selectedOptionValue) {
					resultQuiz.lastElementChild.style.color = "green"
					this.quizzesDone.set(difficulty, this.quizzesDone.get(difficulty) + 1)

				} else {
					resultQuiz.lastElementChild.style.color = "gray"
				}
				return
			}
			if (optionValue == selectedOptionValue && selectedOptionValue != correctOptionValue) {
				resultQuiz.lastElementChild.style.color = "red"
			}
		})
		let answerID = await this.postAnswer(this.quizID, this.quiz.getAttribute("selected-id"))
		console.log(answerID)
	}

	async postAnswer(quizID, optionID) {
		const response = await fetch("/user/quiz/answer", {
			method: "POST",
			body: JSON.stringify({
				quiz_id: quizID,
				trivia_id: this.quiz.getAttribute("trivia-id"),
				option_id: optionID
			}),
			headers: {
				"Content-type": "application/json"
			}
		})
		const answerID = await response.json()
		return answerID
	}

	nextQuiz() {
		if (!this.quiz.getAttribute("selected")) {
			return false
		}
		this.storeResults()
		if (this.index == this.quizzes.length - 1) {
			this.finishQuiz()
			return
		}
		this.index += 1
		this.quiz = this.quizzes[this.index]
		this.quizContainer.replaceChildren(this.quiz)
		this.setUpOptions()
	}

	finishQuiz() {
		let completedIn = (Date.now() - this.startedAt) / 1000
		let score = 0
		let coeff = 1
		for (let v of this.quizzesDone.values()) {
			score += v * coeff
			coeff += 1
		}
		let totalScore = score / completedIn * 1000
		totalScore = Math.round(totalScore)

		this.quizContainer.innerHTML = ""
		this.nextBtn.remove()
		let stats = document.createElement("div")
		stats.className = "stats"
		stats.insertAdjacentHTML("beforeend", ` <h3> Completed in: ${completedIn} seconds</h3> `)
		stats.insertAdjacentHTML("beforeend", `<h3>Score: ${totalScore}</h3>`)

		this.results.appendChild(stats)
		let finishButton = this.results.querySelector(".done-btn")
		this.results.appendChild(finishButton)
		this.results.classList.add("show")
		this.results.hidden = false

		if (this.quizID) {
			fetch("/user/quiz/score", {
				method: "POST",
				body: JSON.stringify({
					completion_time: completedIn * 1000,
					hard_quizzes_done: this.quizzesDone.get("hard"),
					medium_quizzes_done: this.quizzesDone.get("medium"),
					easy_quizzes_done: this.quizzesDone.get("hard"),
					total_score: totalScore,
					quiz_id: this.quizID
				}),
				headers: {
					"Content-type": "application/json"
				}
			}).then(response => response.text()).then(result => console.log(result))
		}
	}
}

new Quizgo()
