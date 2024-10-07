package quizgo

type Trivia struct {
	Type          string   `json:"type"`
	Difficulty    string   `json:"difficulty"`
	Category      string   `json:"category"`
	Question      string   `json:"question"`
	CorrectOption string   `json:"correct_answer"`
	Options       []string `json:"incorrect_answers"`
}
