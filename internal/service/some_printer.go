package service

import (
	"dilema/internal/action"
	"log"
)

type somePrinterWithoutParams struct {
	text string
}

func NewSomePrinterWithoutParams() action.SomePrinter {
	return &somePrinterWithoutParams{text: "Hello world!"}
}

func (sp *somePrinterWithoutParams) PrintSome() {
	log.Println(sp.text)
}
