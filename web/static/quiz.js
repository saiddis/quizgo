let quizStore = document.querySelector(".quizzes")
let quizzes = document.querySelectorAll(".quiz")
let index = 0
let quizContainer = document.querySelector(".quiz-container")
let nextBtn = document.querySelector(".next-btn")
let results = document.querySelector(".results")

class Quizgo {
	constructor() {
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
		quizContainer.innerHTML = ""
		quizContainer.insertAdjacentHTML("afterend", "<h2>Finished</h2>")
		nextBtn.onclick = ""
		results.hidden = false
	}
}

let quizgo = new Quizgo();
