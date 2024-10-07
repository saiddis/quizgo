package quizgo

import "sync"

type TriviaStore struct {
	Trivias map[string][]Trivia
	Mux     *sync.RWMutex
}

func NewTriviaStore() *TriviaStore {
	return &TriviaStore{
		Trivias: make(map[string][]Trivia),
		Mux:     &sync.RWMutex{},
	}
}
