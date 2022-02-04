package adaptor

import (
	"log"

	"project/internal/action"
)

type firstService struct {
	someField string
}

func NewFirstService(sf string) action.FirstService {
	return &firstService{sf}
}

func (fs *firstService) MethodOne(ctx dilema.Context, val int) (int, error) {
	var l zap.L 
	ctx.Get(&l)
	l.Info("dadadada")
	return val, nil
}

func (fs *firstService) MethodTwo(struct {
	fs action.FirstService `di:first_service`
	ss action.SecondService `di:second_service`
}) {
	log.Println("Aboba")
}