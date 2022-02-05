package main

import (
	"dilema"
	"dilema/example/internal/service"
)

func main() {
	// {
	// 	diFirst := dilema.Init()
	// 	err := diFirst.RegisterTemporal("action", service.NewSomeActionByWithParams)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	sum := diFirst.Get(new(action.SomeAction), 1, 3).(action.SomeAction).Sum()
	// 	log.Println(sum)
	// }

	// {
	// 	diThird := dilema.Init()
	// 	err := diThird.RegisterTemporal(service.NewSomeActionByWithParams)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	sum := diThird.Get(new(action.SomeAction), 1, 3).(action.SomeAction).Sum()
	// 	log.Println(sum)

	// }
	{
		diSecond := dilema.Init()
		err := diSecond.RegisterFew(
			map[string]interface{}{
				"action": service.NewSomeActionWithoutParams, 
				"printer": service.NewSomePrinterWithoutParams,
			},
		)
		if err != nil {
			panic(err)
		}
	}
}
