package dilema

import (
	"github.com/kjushka/dilema/dilerr"
	"reflect"
)

func checkProvidedTypeIsCreator(provided interface{}) (reflect.Type, reflect.Value, error) {
	v := reflect.ValueOf(provided)
	t := reflect.TypeOf(provided)

	if t.Kind() != reflect.Func {
		return nil, reflect.Value{},
			dilerr.NewTypeError("expected provided type is func")
	}
	if t.NumOut() < 1 && t.NumOut() > 2 &&
		t.Out(0).Kind() != reflect.Interface &&
		t.Out(0).Kind() != reflect.Ptr &&
		t.Out(0).Kind() != reflect.Struct {
		return t, v,
			dilerr.NewTypeError("expected provided return one or two interface typed value")
	}
	if t.NumOut() == 2 && t.Out(1).String() != "error" {
		return t, v,
			dilerr.NewTypeError("expected provided return second value interface as error")
	}

	return t, v, nil
}

func checkIsError(possibleError reflect.Value) (error, bool) {
	if possibleError.CanInterface() {
		if err, ok := possibleError.Interface().(error); ok {
			return err, ok
		}
	}

	return nil, false
}

func (di *dicon) createInStruct(sType reflect.Type) (reflect.Value, bool) {
	newValue := reflect.New(sType)
	elem := newValue.Elem()

	for i := 0; i < sType.NumField(); i++ {
		if !elem.Field(i).CanSet() {
			return newValue, false
		}
		if alias := sType.Field(i).Tag.Get("dilema"); alias != "" {
			container, ok := di.getSingleToneByAlias(alias)
			if ok {
				elem.Field(i).Set(container)
				continue
			}
		} else {
			fieldType := sType.Field(i).Type
			container, ok := di.getSingleToneByType(fieldType)
			if ok {
				elem.Field(i).Set(container)
				continue
			}
		}
	}

	return newValue, true
}

func processValue(val reflect.Value, container interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	vCont := reflect.ValueOf(container)
	tCont := vCont.Type()
	if tCont.Kind() != reflect.Ptr {
		err = dilerr.NewProcessError("expected ptr values")
		return
	}
	elem := vCont.Elem()
	if !elem.CanSet() {
		err = dilerr.NewProcessError("agruments can't be setted")
		return
	}
	elem.Set(val)
	return
}
