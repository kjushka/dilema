package dilema

import (
	"dilema/dilema/dilerr"
	"log"
	"reflect"
)

func checkProvidedTypeIsCreator(provided interface{}) (reflect.Type, reflect.Value, error) {
	v := reflect.ValueOf(provided)
	t := reflect.TypeOf(provided)

	if t.Kind() != reflect.Func {
		log.Println(t.Kind())
		return nil, reflect.Value{},
			dilerr.NewTypeError("expected provided type is func")
	}
	if t.NumOut() != 1 && t.Out(0).Kind() != reflect.Interface {
		return t, v,
			dilerr.NewTypeError("expected provided return only one interface typed value")
	}
	return t, v, nil
}

func checkProvidedTypeIsUnion(provided interface{}) (reflect.Type, reflect.Value, error) {
	v := reflect.ValueOf(provided).Elem()
	t := reflect.TypeOf(provided).Elem()

	if t.Kind() != reflect.Interface {
		return t, v, dilerr.NewTypeError("expected provided type is interface")
	}

	return t, v, nil
}
