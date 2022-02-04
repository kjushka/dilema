package main

import (
	"log"
	"project/internal/adaptor"

	"github.com/kjushka/dilema"
	"github.com/postgres/postgres"
	"github.com/uber/zap"
)

var Di dilema.Dicon

func init() {
    Di = dilema.New()
}

func main() {
    config()
    Di.Recover(run)
}

func config() {
    //init database
    db := postgres.New("127.0.0.1:5432")
    Di.Register(db, "postgres")

    //add smth to ctx
    l := zap.New()
    Di.AddToCtx(l)

    Di.RegisterSingleTone("first_service", NewFirstService, "hello_string")
    Di.RegisterTemporal("second_service", NewSecondService)
}

func run() {
    var (
        err error
        val int
        fs FirstService
    )

    Di.Run(fs.MethodTwo)

    Di.Get("first_service").Run(fs.FirstMethod, "").PanicIfError().Process(&fs)
    Di.Run(fs.MethodOne, 123).Process(&val, &err)
    if err != nil {
        panic(err)
    }
    log.Println(val)
}