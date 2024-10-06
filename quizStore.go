package quizgo

import "sync"

type QuizStore struct {
	Quizzes map[string][]Quiz
	Mux     *sync.RWMutex
}

func NewQuizStore() *QuizStore {
	return &QuizStore{
		Quizzes: make(map[string][]Quiz),
		Mux:     &sync.RWMutex{},
	}
}
