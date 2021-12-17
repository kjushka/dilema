package main

import (
	"dilema/dilema"
	"dilema/internal/action"
	"dilema/internal/service"
	"log"
)

func main() {
	diFirst := dilema.Init()
	err := diFirst.ProvideSingleTone(service.NewSomeActionByWithParams, 1, 3)
	if err != nil {
		panic(err)
	}

	sum := diFirst.Get(new(action.SomeAction)).(action.SomeAction).Sum()
	log.Println(sum)

	diSecond := dilema.Init()
	err = diSecond.ProvideAll(service.NewSomeActionWithoutParams, service.NewSomePrinterWithoutParams)
	log.Println("SUM:", diSecond.Get(new(action.SomeAction)).(action.SomeAction).Sum())
	diSecond.Get(new(action.SomePrinter)).(action.SomePrinter).PrintSome()
}
