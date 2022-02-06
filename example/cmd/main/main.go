package main

import (
	"dilema"
	"dilema/example/internal/action"
	"dilema/example/internal/service"
	"log"
)

func main() {
	{
		diFirst := dilema.Init()
		err := diFirst.RegisterTemporal("action", service.NewSomeActionByWithParams)
		if err != nil {
			panic(err)
		}

		container, err := diFirst.GetTemporal("action", 1, 3)
		if err != nil {
			panic(err)
		}
		sum := container.(action.SomeAction).Sum()
		log.Println(sum)
	}
	{
		diSecond := dilema.Init()
		diSecond.MustRegisterTemporal("action", service.NewSomeActionByWithParams)
		sum := diSecond.MustGetTemporal("action", 1, 3).(action.SomeAction).Sum()
		log.Println(sum)
	}
	{
		diThird := dilema.Init()
		err := diThird.RegisterFew(
			map[string]interface{}{
				"action":  service.NewSomeActionWithoutParams,
				"printer": service.NewSomePrinterWithoutParams,
			},
		)
		if err != nil {
			panic(err)
		}
		sum := diThird.MustGetSingletone("action").(action.SomeAction).Sum()
		log.Println(sum)
		printer, err := diThird.GetSingletone("printer")
		if err != nil {
			panic(err)
		}
		printer.(action.SomePrinter).PrintSome()

		diThird.MustRun(someFunc)
	}
}

func someFunc(diStruct *struct {
	Action  action.SomeAction  `dilema:"action"`
	Printer action.SomePrinter `dilema:"printer"`
}) {
	log.Println("inside someFunc:", diStruct.Action.Sum())
	diStruct.Printer.PrintSome()
}
