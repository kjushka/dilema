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
	if reflect.TypeOf(destroyer).Kind() != reflect.Func {
		return dilerr.NewTypeError("Unexpected destroyer type")
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

func (di *dicon) checkCreationResults(creationResults []reflect.Value) (int, error) {
	var destroyerIndex int
	if len(creationResults) > 1 {
		err, errIndex := checkHasError(creationResults)
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

func checkHasError(creationResults []reflect.Value) (error, int) {
	for i, result := range creationResults {
		if result.CanInterface() {
			if err, ok := result.Interface().(error); ok {
				return err, i
			}
		}
	}

	return nil, -1
}