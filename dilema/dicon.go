package dilema

import (
	"dilema/dilema/dilerr"
	"reflect"
	"sync"
)

type dicon struct {
	temps       map[reflect.Type]Constructor
	singleTones map[reflect.Type]reflect.Value

	functions map[reflect.Method]actionToType
	cache     map[reflect.Type]reflect.Value

	mutex *sync.Mutex
}

type actionToType struct {
	actionType  reflect.Type
	serviceType serviceType
}

// ProvideTemporal provides new service, which will be initialized when
// you call Get method and be destroyed with GC after work will be done
func (di *dicon) ProvideTemporal(serviceInit interface{}) error {
	t, v, err := checkProvidedTypeIsCreator(serviceInit)
	if err != nil {
		return err
	}

	di.temps[t.Out(0)] = &constructor{creator: v}
	di.registerFunctions(t.Out(0), temporalType)
	return nil
}

// ProvideSingleTone provides new singletone - constant service, which is being created only
// one time during all time that program works. It's being initialized immediately
// ProvideSingleTone called. That's why if it's necessary some arguments can be attached
func (di *dicon) ProvideSingleTone(serviceInit interface{}, args ...interface{}) error {
	t, v, err := checkProvidedTypeIsCreator(serviceInit)

	argsIndex := 0
	singleTone, err := di.createService(v, &argsIndex, args...)
	if err != nil {
		return err
	}

	di.singleTones[t.Out(0)] = singleTone
	di.registerFunctions(t.Out(0), singleToneType)
	return nil
}

// ProvideAll provides some amount of services, which can be initialized without extra arguments.
// It's equal calling ProvideSingleTone for every several service
func (di *dicon) ProvideAll(servicesInit ...interface{}) error {
	servicesMap := make(map[reflect.Type]reflect.Value)
	for _, serviceInit := range servicesInit {
		t, v, err := checkProvidedTypeIsCreator(serviceInit)
		if err != nil {
			return err
		}
		t = t.Out(0)
		servicesMap[t] = v
	}

	type tv struct {
		t reflect.Type
		v reflect.Value
	}

	services := make([]tv, 0)
	for t, v := range servicesMap {
		service, err := di.createService(v, nil)
		if err != nil {
			return err
		}
		services = append(services, tv{t, service})
	}

	for _, service := range services {
		di.singleTones[service.t] = service.v
		di.registerFunctions(service.t, singleToneType)
	}

	return nil
}

// createService creates instance of service, which interface return from provided func
func (di *dicon) createService(v reflect.Value, argsIndex *int, args ...interface{}) (reflect.Value, error) {
	t := v.Type()
	ins := make([]reflect.Value, 0)
	for i := 0; i < t.NumIn(); i += 1 {
		in, ok := di.checkInDiconServices(t, i, argsIndex, args...)
		if !ok {
			if argsIndex != nil &&
				len(args)-1 >= *argsIndex &&
				reflect.TypeOf(args[*argsIndex]) == t.In(i) {
				in = reflect.ValueOf(args[*argsIndex])
				*(argsIndex) += 1
			} else {
				return reflect.Value{},
					dilerr.NewCreationError("there are no necessary arguments to create a service")
			}
		}
		ins = append(ins, in)
	}

	return v.Call(ins)[0], nil
}

// checkInDiconServices checks containing of necessary typed service in dicon services
func (di *dicon) checkInDiconServices(
	t reflect.Type,
	i int,
	argsIndex *int,
	args ...interface{},
) (reflect.Value, bool) {
	paramT := t.In(i)
	temp, ok := di.temps[paramT]
	if ok {
		service, err := di.createService(temp.Creator(), argsIndex, args...)
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

func (di *dicon) registerFunctions(action reflect.Type, sType serviceType) {
	for i := 0; i < action.NumMethod(); i += 1 {
		method := action.Method(i)
		fn, ok := di.functions[method]
		if ok && fn.actionType.NumMethod() > action.NumMethod() {
			continue
		}

		method.Index = 0
		di.mutex.Lock()
		di.functions[method] = actionToType{
			actionType:  action,
			serviceType: sType,
		}
		di.mutex.Unlock()
	}
}

// Get return services typed with some interface or construct and return service, if it is temporal.
func (di *dicon) Get(serviceAction interface{}, args ...interface{}) interface{} {
	t := reflect.TypeOf(serviceAction).Elem()

	singleTone, ok := di.singleTones[t]
	if ok {
		return singleTone.Convert(t).Interface()
	}

	tempConstructor, ok := di.temps[t]
	if ok {
		argsIndex := 0
		temp, err := di.createService(tempConstructor.Creator(), &argsIndex, args...)
		if err != nil {
			return nil
		}
		return temp.Convert(t).Interface()
	}

	return nil
}

/*// GetFromUnion returnStruct which is implementation of provided interface.
// Provided interface must contain fields typed with provided earlier services
// Now it must contain only singleTones
func (di *dicon) GetFromUnion(union interface{}, args ...interface{}) interface{} {
	t, _, err := checkProvidedTypeIsUnion(union)
	if err != nil {
		return nil
	}

	fields := make([]reflect.StructField, 0)
	fieldsValues := make(map[reflect.Type]reflect.Value)
	argsIndex := 0
	for i := 0; i < t.NumMethod(); i++ {
		t, sType, ok := di.findServiceByMethod(t.Method(i))
		if !ok {
			return nil
		}

		var value reflect.Value
		switch sType {
		case singleToneType:
			value = di.singleTones[t]
		case temporalType:
			constr := di.temps[t]
			value, err = di.createService(constr.Creator(), &argsIndex, args...)
			if err != nil {
				return nil
			}
		}

		field := reflect.StructField{
			Name:      t.Name(),
			Type:      t,
			Anonymous: false,
		}
		fields = append(fields, field)
		fieldsValues[t] = value
	}

	structPtr := reflect.New(reflect.StructOf(fields))
	structValue := structPtr.Elem()

	for i := 0; i < structValue.NumField(); i += 1 {
		field := structValue.Field(i)
		value := fieldsValues[field.Type()]
		structValue.Field(i).Set(value)
	}

	return structValue.Interface()
}*/

func (di *dicon) findServiceByMethod(method reflect.Method) (reflect.Type, serviceType, bool) {
	method.Index = 0
	di.mutex.Lock()
	action, ok := di.functions[method]
	di.mutex.Unlock()
	return action.actionType, action.serviceType, ok
}

func (di *dicon) addToCache(t reflect.Type, v reflect.Value) {
	di.mutex.Lock()
	di.cache[t] = v
	di.mutex.Unlock()
}
