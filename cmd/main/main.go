package main

import (
	"dilema/dilema"
	"dilema/internal/action"
	"dilema/internal/service"
	"log"
)

func main() {
	di := dilema.Init()
	err := di.ProvideSingleTone(service.NewSomeAction, 1, 3)
	if err != nil {
		panic(err)
	}

	sum := di.Get(new(action.SomeAction)).(action.SomeAction).Sum()
	log.Println(sum)
}