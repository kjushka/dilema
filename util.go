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
	if t.NumOut() < 1 && t.Out(0).Kind() != reflect.Interface {
		return t, v,
			dilerr.NewTypeError("expected provided return only one interface typed value")
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

func checkProvidedTypeIsCorrectStruct(provided interface{}) (reflect.Type, reflect.Value, error) {
	v := reflect.ValueOf(provided).Elem()
	t := reflect.TypeOf(provided).Elem()

	if t.Kind() != reflect.Struct {
		return t, v, dilerr.NewTypeError("expected provided type is struct")
	}

	for i := 0; i < t.NumField(); i += 1 {
		if t.Field(i).Type.Kind() != reflect.Interface || t.Field(i).Tag.Get("di") == "" {
			return t, v, dilerr.NewTypeError("expected all fields are interfaces and fields have tags 'di'")
		}
	}

	return t, v, nil
}

func (di *dicon) checkCreationResults(creationResults []reflect.Value) (destroyerIndex int, err error) {
	if len(creationResults) > 1 {
		errIndex, err := checkHasError(creationResults)
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

func checkHasError(creationResults []reflect.Value) (int, error) {
	for i, result := range creationResults {
		if result.CanInterface() {
			if err, ok := result.Interface().(error); ok {
				return i, err
			}
		}
	}

	return -1, nil
}

func (di *dicon) createCorrectInStruct(sType reflect.Type, args ...interface{}) (reflect.Value, bool) {
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
			constuctor, ok := di.getTemporalByAlias(alias)
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
			constuctor, ok := di.getTemporalByType(fieldType)
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
