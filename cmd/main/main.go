package main

import (
	"dilema/dilema"
	"dilema/internal/service"
)

func main() {
	/*{
		diFirst := dilema.Init()
		err := diFirst.ProvideTemporal(service.NewSomeActionByWithParams)
		if err != nil {
			panic(err)
		}

		sum := diFirst.Get(new(action.SomeAction), 1, 3).(action.SomeAction).Sum()
		log.Println(sum)
	}

	{
		diThird := dilema.Init()
		err := diThird.ProvideTemporal(service.NewSomeActionByWithParams)
		if err != nil {
			panic(err)
		}
		sum := diThird.Get(new(action.SomeAction), 1, 3).(action.SomeAction).Sum()
		log.Println(sum)

	}
	*/
	{
		diSecond := dilema.Init()
		err := diSecond.ProvideAll(service.NewSomeActionWithoutParams, service.NewSomePrinterWithoutParams)
		if err != nil {
			panic(err)
		}
	}
}
