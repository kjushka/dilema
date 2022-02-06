package main

import (
	"dilema"
	"dilema/example/internal/action"
	"dilema/example/internal/service"
	"fmt"
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

		res, err := diThird.Recover(someFunc, 1)
		if err != nil {
			log.Println(err.Error())
		}
		var (
			val int
		)
		res.MustProcess(&val, &err)
		log.Println(val, err)
	}
}

func someFunc(diStruct *struct {
	Action  action.SomeAction  `dilema:"action"`
	Printer action.SomePrinter `dilema:"printer"`
}, num int) (val int, err error) {
	val, err = 666, fmt.Errorf("test error")
	log.Println("inside someFunc:", "num:", num)
	log.Println("inside someFunc:", "sum:", diStruct.Action.Sum())
	diStruct.Printer.PrintSome()
	return
}
