package dilema

import (
	"dilema/dilema/dilerr"
	"log"
	"reflect"
)

func checkProvidedTypeIsFunc(provided interface{}) (reflect.Type, reflect.Value, error) {
	v := reflect.ValueOf(provided)
	t := reflect.TypeOf(provided)

	if t.Kind() != reflect.Func {
		log.Println(t.Kind())
		return nil, reflect.Value{},
			dilerr.NewTypeError("unexpected serviceInit type")
	}
	if t.NumOut() != 1 && t.Out(0).Kind() != reflect.Interface {
		return t, v,
			dilerr.NewTypeError("expected serviceInit return only one interface typed value")
	}
	return t, v, nil
}
