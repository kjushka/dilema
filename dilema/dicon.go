package dilema

import (
	"dilema/dilema/dilerr"
	"log"
	"reflect"
)

type dicon struct {
	temps       map[reflect.Type]Constructor
	singleTones map[reflect.Type]reflect.Value
}

// ProvideTemp provides new service, which will be initialized when
// you call Get method and be destroyed with GC after work will be done
func (di *dicon) ProvideTemp(serviceInit interface{}) error {
	t, v, err := checkProvidedTypeIsFunc(serviceInit)
	if err != nil {
		return err
	}

	di.temps[t.Out(0)] = &constructor{creator: v}
	return nil
}

// ProvideSingleTone provides new singletone - constant service, which is being created only
// one time during all time that program works. It's being initialized immediately
// ProvideSingleTone called. That's why if it's necessary some arguments can be attached
func (di *dicon) ProvideSingleTone(serviceInit interface{}, args ...interface{}) error {
	t, v, err := checkProvidedTypeIsFunc(serviceInit)

	singleTone, err := di.createService(v, args...)
	if err != nil {
		return err
	}

	di.singleTones[t.Out(0)] = singleTone
	return nil
}

// ProvideAll provides some amount of services, which can be initialized without extra arguments.
// It's equal calling ProvideSingleTone for every several service
func (di *dicon) ProvideAll(servicesInit ...interface{}) error {
	servicesMap := make(map[reflect.Type]reflect.Value)
	for _, serviceInit := range servicesInit {
		t, v, err := checkProvidedTypeIsFunc(serviceInit)
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
		service, err := di.createService(v)
		if err != nil {
			return err
		}
		services = append(services, tv{t, service})
	}

	for _, service := range services {
		di.singleTones[service.t] = service.v
	}

	log.Println(di.singleTones)
	return nil
}

// createService creates instance of service, which interface return from provided func
func (di *dicon) createService(v reflect.Value, args ...interface{}) (reflect.Value, error) {
	t := v.Type()
	ins := make([]reflect.Value, 0)
	argsIndex := 0
	for i := 0; i < t.NumIn(); i++ {
		in, ok := di.checkInDiconServices(t, i)
		if !ok {
			if len(args)-1 >= argsIndex && reflect.TypeOf(args[argsIndex]) == t.In(i) {
				in = reflect.ValueOf(args[argsIndex])
				argsIndex += 1
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
func (di *dicon) checkInDiconServices(t reflect.Type, i int) (reflect.Value, bool) {
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

// Get return services typed with some interface or construct and return service, if it is temporal.
func (di *dicon) Get(serviceAction interface{}, args ...interface{}) interface{} {
	t := reflect.TypeOf(serviceAction).Elem()

	singleTone, ok := di.singleTones[t]
	if ok {
		return singleTone.Convert(t).Interface()
	}

	tempConstructor, ok := di.temps[t]
	if ok {
		temp, err := di.createService(tempConstructor.Creator(), args...)
		if err != nil {
			return nil
		}
		return temp.Convert(t).Interface()
	}

	return nil
}
