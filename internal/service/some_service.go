package service

import "dilema/internal/action"

type someService struct {
	a, b int
}

func NewSomeAction(a, b int) action.SomeAction {
	return &someService{a, b}
}

func (ss *someService) Sum() int {
	return ss.a + ss.b
}
