let quizStore = document.querySelector(".quizzes")
let quizzes = document.querySelectorAll(".quiz")
let index = 0
let quizContainer = document.querySelector(".quiz-container")
let nextBtn = document.querySelector(".next-btn")
let results = document.querySelector(".results-container")
let quizID = quizContainer.getAttribute("quiz_id")

class Quizgo {
	constructor() {
		this.startedAt = Date.now()
		this.quizzesDone = new Map([
			["easy", 0],
			["medium", 0],
			["hard", 0]
		]);

		this.quiz = quizzes[index]
		this.setUpOptions()

		quizContainer.append(this.quiz)
		nextBtn.onclick = () => {
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
			}
		})
	}

	storeResults() {
		let correctOptionValue = this.quiz.querySelector(".correct").getAttribute("value")
		let selectedOptionValue = this.quiz.getAttribute("selected")
		let difficulty = this.quiz.getAttribute("difficulty")

		let stats = document.createElement("h3")
		stats.innerText = `
Quiz #${index + 1}: ${this.quiz.querySelector(".question").innerHTML}
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

			results.append(resultQuiz)

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
	}

	nextQuiz() {
		if (!this.quiz.getAttribute("selected")) {
			return false
		}
		this.storeResults()
		if (index == quizzes.length - 1) {
			this.finishQuiz()
			return
		}
		index += 1
		this.quiz = quizzes[index]
		quizContainer.replaceChildren(this.quiz)
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

		quizContainer.innerHTML = ""
		nextBtn.remove()
		let stats = document.createElement("div")
		stats.className = "stats"
		stats.insertAdjacentHTML("beforeend", ` <h3> Completed in: ${completedIn} seconds</h3> `)
		stats.insertAdjacentHTML("beforeend", `<h3>Score: ${totalScore}</h3>`)

		results.appendChild(stats)
		let finishButton = results.querySelector(".done-btn")
		results.appendChild(finishButton)
		results.classList.add("show")
		results.hidden = false

		if (quizID) {
			fetch("/user/quiz/score", {
				method: "POST",
				body: JSON.stringify({
					completion_time: completedIn * 1000,
					hard_quizzes_done: this.quizzesDone.get("hard"),
					medium_quizzes_done: this.quizzesDone.get("medium"),
					easy_quizzes_done: this.quizzesDone.get("hard"),
					total_score: totalScore,
					quiz_id: quizID
				}),
				headers: {
					"Content-type": "application/json"
				}
			}).then(response => response.text()).then(result => console.log(result))
		}
	}
}

let quizgo = new Quizgo();
