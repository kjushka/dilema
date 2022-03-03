package main

import (
	"dilema"
	"dilema/example/internal/action"
	"dilema/example/internal/service"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	{
		diFirst := dilema.Init()
		err := diFirst.RegisterTemporary("action", service.NewSomeActionByWithParams)
		if err != nil {
			panic(err)
		}

		container, err := diFirst.GetTemporary("action", 1, 3)
		if err != nil {
			panic(err)
		}
		sum := container.(action.SomeAction).Sum()
		log.Println(sum)
	}
	{
		diSecond := dilema.Init()
		diSecond.MustRegisterTemporary("action", service.NewSomeActionByWithParams)
		sum := diSecond.MustGetTemporary("action", 1, 3).(action.SomeAction).Sum()
		log.Println(sum)
	}
	{
		diThird := dilema.Init()
		diThird.MustRegisterSingletone("action", service.NewSomeActionWithoutParams)
		diThird.MustRegisterSingletone("printer", service.NewSomePrinterWithoutParams)

		sum := diThird.MustGetSingletone("action").(action.SomeAction).Sum()
		log.Println(sum)
		printer, err := diThird.GetSingletone("printer")
		if err != nil {
			panic(err)
		}
		printer.(action.SomePrinter).PrintSome()

		var file io.Reader
		file, err = os.Open("./example/cmd/main/text.txt")
		if err != nil {
			panic(err)
		}
		res, err := diThird.RecoverAndClean(someFunc, 1, file)
		if err != nil {
			log.Println(err.Error())
			panic(err)
		}
		var (
			val int
		)
		res.MustProcess(&val, &err)
		log.Println(val, err)
	}
}

func someFunc(num int, diStruct *struct {
	Action  action.SomeAction  `dilema:"action"`
	Printer action.SomePrinter `dilema:"printer"`
}, r io.Reader) (val int, err error) {
	val = 666
	log.Println("inside someFunc:", "num:", num)
	log.Println("inside someFunc:", "sum:", diStruct.Action.Sum())
	diStruct.Printer.PrintSome()
	var buf []byte
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		panic(err)
	}
	err = fmt.Errorf("test error")
	log.Println("readed bytes:", n)
	log.Println("readed string:", string(buf))
	return
}
