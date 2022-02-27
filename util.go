package dilema

import (
	"dilema/dilerr"
	"reflect"
)

func checkProvidedTypeIsCreator(provided interface{}) (reflect.Type, reflect.Value, error) {
	v := reflect.ValueOf(provided)
	t := reflect.TypeOf(provided)

	if t.Kind() != reflect.Func {
		return nil, reflect.Value{},
			dilerr.NewTypeError("expected provided type is func")
	}
	if t.NumOut() < 1 && t.NumOut() > 2 && t.Out(0).Kind() != reflect.Interface {
		return t, v,
			dilerr.NewTypeError("expected provided return one or two interface typed value")
	}
	if t.NumOut() == 2 && t.Out(1).String() != "error" {
		return t, v,
			dilerr.NewTypeError("expected provided return second value interface as error")
	}

	return t, v, nil
}

func checkIsDestroyer(destroyer interface{}) error {
	dType := reflect.TypeOf(destroyer)
	if dType.Kind() != reflect.Func {
		return dilerr.NewTypeError("Unexpected destroyer type")
	}
	if dType.NumIn() != 0 {
		return dilerr.NewTypeError("Destroyer must have no in arguments")
	}
	return nil
}

func (di *dicon) checkCreationResults(creationResults []reflect.Value) (destroyerIndex int, err error) {
	if len(creationResults) > 1 {
		errIndex, err := checkIsError(creationResults)
		if errIndex != -1 && err != nil {
			return -1, err
		}
		if len(creationResults) == 2 && errIndex == -1 {
			err = checkIsDestroyer(creationResults[1])
			if err != nil {
				return -1, err
			}
			destroyerIndex = 1
		} else if len(creationResults) == 3 && errIndex != -1 {
			dIndex := 3 - errIndex
			err = checkIsDestroyer(creationResults[dIndex])
			if err != nil {
				return -1, err
			}
			destroyerIndex = dIndex
		} else {
			return -1, dilerr.NewTypeError("Creator have too much outs")
		}
	}
	return destroyerIndex, nil
}

func checkIsError(possibleError reflect.Value) (error, bool) {
	if possibleError.CanInterface() {
		if err, ok := possibleError.Interface().(error); ok {
			return err, ok
		}
	}

	return nil, false
}

func (di *dicon) createInStruct(sType reflect.Type, args ...interface{}) (reflect.Value, bool) {
	newValue := reflect.New(sType)
	elem := newValue.Elem()

	for i := 0; i < sType.NumField(); i++ {
		if sType.Field(i).Type.Kind() != reflect.Interface {
			return newValue, false
		}
		if !elem.Field(i).CanSet() {
			return newValue, false
		}
		if alias := sType.Field(i).Tag.Get("dilema"); alias != "" {
			container, ok := di.getSingleToneByAlias(alias)
			if ok {
				elem.Field(i).Set(container)
				continue
			}
			constuctor, ok := di.getTemporaryByAlias(alias)
			if ok {
				argsIndex := 0
				creationResults, err := di.createService(constuctor, &argsIndex, args...)
				if err != nil {
					return newValue, false
				}
				errIndex, err := checkHasError(creationResults)
				if errIndex != -1 && err != nil {
					return newValue, false
				}

				elem.Field(i).Set(creationResults[0])
				continue
			}
		} else {
			fieldType := sType.Field(i).Type
			container, ok := di.getSingleToneByType(fieldType)
			if ok {
				elem.Field(i).Set(container)
				continue
			}
			constuctor, ok := di.getTemporaryByType(fieldType)
			if ok {
				argsIndex := 0
				creationResults, err := di.createService(constuctor, &argsIndex, args...)
				if err != nil {
					return newValue, false
				}
				errIndex, err := checkHasError(creationResults)
				if errIndex != -1 && err != nil {
					return newValue, false
				}

				elem.Field(i).Set(creationResults[0])
				continue
			}
		}
	}

	return newValue, true
}

func processValue(val reflect.Value, container interface{}) error {
	vCont := reflect.ValueOf(container)
	tCont := vCont.Type()
	elem := vCont.Elem()
	if tCont.Kind() != reflect.Ptr {
		return dilerr.NewTypeError("expected ptr values")
	}
	if !elem.CanSet() {
		return dilerr.NewTypeError("agruments can't be setted")
	}
	elem.Set(val)
	return nil
}
