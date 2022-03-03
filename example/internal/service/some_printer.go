package service

import (
	"log"
)

type SomePrinterWithoutParams struct {
	text string
}

func NewSomePrinterWithoutParams() *SomePrinterWithoutParams {
	return &SomePrinterWithoutParams{text: "Hello world!"}
}

func (sp *SomePrinterWithoutParams) PrintSome() {
	log.Println(sp.text)
}
