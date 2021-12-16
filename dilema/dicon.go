package dilema

import (
	"errors"
	"reflect"
)

type dicon struct {
	temps       map[reflect.Type]Constructor
	singleTones map[reflect.Type]reflect.Value
}

func (di *dicon) ProvideTemp(serviceInit interface{}) error {
	if reflect.TypeOf(serviceInit).Kind() != reflect.Func {
		return errors.New("unexpected serviceInit type")
	}

	v := reflect.ValueOf(serviceInit)
	t := reflect.TypeOf(serviceInit)
	if t.NumOut() < 1 {
		return errors.New("expected only one return value")
	}

	if t.Out(0).Kind() != reflect.Interface {
		return errors.New("expected returning of service interface")
	}

	di.temps[t.Out(0)] = &constructor{creator: v}
	return nil
}

func (di *dicon) ProvideSingleTone(serviceInit interface{}, args ...interface{}) error {
	if reflect.TypeOf(serviceInit).Kind() != reflect.Func {
		return errors.New("unexpected serviceInit type")
	}

	v := reflect.ValueOf(serviceInit)
	t := reflect.TypeOf(serviceInit)
	if t.NumOut() != 1 {
		return errors.New("expected only one return value")
	}

	if t.Out(0).Kind() != reflect.Interface {
		return errors.New("expected returning of service interface")
	}

	singleTone, err := di.createService(v, args...)
	if err != nil {
		return err
	}

	di.singleTones[t.Out(0)] = singleTone
	return nil
}

func (di *dicon) createService(v reflect.Value, args... interface{}) (reflect.Value, error) {
	t := v.Type()
	ins := make([]reflect.Value, 0)
	argsIndex := 0
	for i := 0; i < t.NumIn(); i++ {
		in, ok := di.checkInParam(t, i)
		if !ok {
			if reflect.TypeOf(args[argsIndex]) == t.In(i) {
				in = reflect.ValueOf(args[argsIndex])
				argsIndex += 1
			} else {
				return reflect.Value{}, errors.New("there are no necessary arguments to create a service")
			}
		}
		ins = append(ins, in)
	}

	return v.Call(ins)[0], nil
}

func (di *dicon) checkInParam(t reflect.Type, i int) (reflect.Value, bool) {
	paramT := t.In(i)
	temp, ok := di.temps[paramT]
	if ok {
		service, err := di.createService(temp.Creator())
		if err != nil {
			return reflect.Value{}, false
		}
		return service, true
	}
	singleTone, ok := di.singleTones[paramT]
	if ok {
		return singleTone, ok
	}
	return reflect.Value{}, ok
}

func (di *dicon) Get(serviceAction interface{}, args... interface{}) interface{} {
	t := reflect.TypeOf(serviceAction).Elem()

	singleTone, ok := di.singleTones[t]
	if ok {
		return singleTone.Convert(t).Interface()//, ok
	}

	tempConstructor, ok := di.temps[t]
	if ok {
		temp, err := di.createService(tempConstructor.Creator(), args...)
		if err != nil {
			return nil//, false
		}
		return temp.Convert(t).Interface()//, ok
	}

	return nil//, false
}