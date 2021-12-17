package service

import (
	"dilema/internal/action"
	"math/rand"
)

type someServiceWithParams struct {
	a, b int
}

func NewSomeActionByWithParams(a, b int) action.SomeAction {
	return &someServiceWithParams{a, b}
}

func (ssp *someServiceWithParams) Sum() int {
	return ssp.a + ssp.b
}

type someServiceWithoutParams struct {
	a, b int
}

func NewSomeActionWithoutParams() action.SomeAction {
	return &someServiceWithoutParams{rand.Int(), rand.Int()}
}

func (ss *someServiceWithoutParams) Sum() int {
	return ss.a + ss.b
}
