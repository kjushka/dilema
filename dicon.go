package dilema

import (
	"context"
	"dilema/dilerr"
	"fmt"
	"reflect"
	"sync"
)

type dicon struct {
	temporalByAlias    map[string]reflect.Value
	temporalByType     map[reflect.Type]reflect.Value
	singletonesByAlias map[string]reflect.Value
	singletonesByType  map[reflect.Type]reflect.Value

	destroyers []reflect.Value
	cache      map[reflect.Type]reflect.Value

	mutex *sync.Mutex
	ctx   context.Context
}

func (di *dicon) RegisterTemporal(alias string, serviceInit interface{}) error {
	return di.registerTemporal(alias, serviceInit)
}

func (di *dicon) MustRegisterTemporal(alias string, serviceInit interface{}) {
	err := di.registerTemporal(alias, serviceInit)
	if err != nil {
		panic(err)
	}
}

// registerTemporal provides new service, which will be initialized when
// you call Get method and be destroyed with GC after work will be done
func (di *dicon) registerTemporal(alias string, serviceInit interface{}) error {
	_, v, err := checkProvidedTypeIsCreator(serviceInit)
	if err != nil {
		return err
	}

	di.temporalByAlias[alias] = v
	return nil
}

func (di *dicon) RegisterSingletone(
	alias string,
	serviceInit interface{},
	args ...interface{},
) error {
	return di.registerSingleTone(alias, serviceInit, args...)
}

func (di *dicon) MustRegisterSingletone(
	alias string,
	serviceInit interface{},
	args ...interface{},
) {
	err := di.registerSingleTone(alias, serviceInit, args...)
	if err != nil {
		panic(err)
	}
}

// registerSingleTone provides new singletone - constant service, which is being created only
// one time during all time that program works. It's being initialized immediately
// ProvideSingleTone called. That's why if it's necessary some arguments can be attached
func (di *dicon) registerSingleTone(
	alias string,
	serviceInit interface{},
	args ...interface{},
) error {
	t, v, err := checkProvidedTypeIsCreator(serviceInit)

	argsIndex := 0
	creationResults, err := di.createService(v, &argsIndex, args...)
	if err != nil {
		return err
	}
	destroyerIndex, err := di.checkCreationResults(creationResults)
	if err != nil {
		return err
	}

	di.singletonesByAlias[alias] = creationResults[0]
	di.singletonesByType[t] = creationResults[0]
	di.destroyers = append(di.destroyers, creationResults[destroyerIndex])

	return nil
}

func (di *dicon) RegisterFew(servicesInit map[string]interface{}, args ...interface{}) error {
	return di.registerFew(servicesInit, args...)
}

func (di *dicon) MustRegisterFew(servicesInit map[string]interface{}, args ...interface{}) {
	err := di.registerFew(servicesInit, args...)
	if err != nil {
		panic(err)
	}
}

// RegisterFew provides some amount of services, which can be initialized without extra arguments.
// It's equal calling ProvideSingleTone for every several service
func (di *dicon) registerFew(servicesInit map[string]interface{}, args ...interface{}) error {
	servicesMap := make(map[string]reflect.Value)
	for alias, serviceInit := range servicesInit {
		_, v, err := checkProvidedTypeIsCreator(serviceInit)
		if err != nil {
			return err
		}
		servicesMap[alias] = v
	}

	type ta struct {
		a string
		t reflect.Type
		v reflect.Value
		d interface{}
	}

	services := make([]ta, 0)
	for a, v := range servicesMap {
		argsIndex := 0
		creationResults, err := di.createService(v, &argsIndex, args...)
		if err != nil {
			return err
		}

		destroyerIndex, err := di.checkCreationResults(creationResults)
		if err != nil {
			return err
		}

		services = append(services, ta{
			a,
			creationResults[0].Type(),
			creationResults[0],
			creationResults[destroyerIndex],
		})
	}

	for _, service := range services {
		di.singletonesByAlias[service.a] = service.v
		di.singletonesByType[service.t] = service.v
		if service.d != nil {
			destroyer := service.d.(reflect.Value)
			di.destroyers = append(di.destroyers, destroyer)
		}
	}

	return nil
}

// createService creates instance of service, which interface return from provided func
func (di *dicon) createService(
	v reflect.Value,
	argsIndex *int,
	args ...interface{},
) ([]reflect.Value, error) {
	t := v.Type()
	ins := make([]reflect.Value, 0)
	for i := 0; i < t.NumIn(); i += 1 {
		in, ok, err := di.checkInDiconServices(t, i, argsIndex, args...)
		if err != nil {
			return nil, err
		}
		if !ok {
			if argsIndex != nil &&
				len(args)-1 >= *argsIndex &&
				reflect.TypeOf(args[*argsIndex]) == t.In(i) {
				in = reflect.ValueOf(args[*argsIndex])
				*(argsIndex) += 1
			} else {
				return nil, dilerr.NewCreationError(
					"there are no necessary arguments to create a service",
				)
			}
		}
		ins = append(ins, in)
	}

	return v.Call(ins), nil
}

// checkInDiconServices checks containing of necessary typed service in dicon services
func (di *dicon) checkInDiconServices(
	t reflect.Type,
	i int,
	argsIndex *int,
	args ...interface{},
) (reflect.Value, bool, error) {
	paramT := t.In(i)
	temp, ok := di.temporalByType[paramT]
	if ok {
		creationResults, err := di.createService(temp, argsIndex, args...)
		if err != nil {
			return reflect.Value{}, false, err
		}
		return creationResults[0], true, nil
	}
	singleTone, ok := di.singletonesByType[paramT]
	if ok {
		return singleTone, ok, nil
	}
	return reflect.Value{}, ok, nil
}

func (di *dicon) GetSingletone(alias string) (interface{}, error) {
	container, err := di.getSingletone(alias)
	if err != nil {
		return nil, err
	}
	return container.Interface(), nil
}

func (di *dicon) MustGetSingletone(alias string) interface{} {
	container, err := di.getSingletone(alias)
	if err != nil {
		panic(err)
	}
	return container.Interface()
}

func (di *dicon) ProcessSingletone(alias string, container interface{}) error {
	c, err := di.getSingletone(alias)
	v := reflect.ValueOf(container)
	if err == nil {
		v.Set(c)
	} else {
		v.Set(reflect.Zero(v.Type()))
	}

	return err
}

func (di *dicon) MustProcessSingletone(alias string, container interface{}) {
	c, err := di.getSingletone(alias)
	v := reflect.ValueOf(container)
	if err != nil {
		panic(err)
	}
	v.Set(c)
}

func (di *dicon) getSingletone(alias string) (reflect.Value, error) {
	singleTone, ok := di.singletonesByAlias[alias]
	if ok {
		return singleTone, nil
	}
	return reflect.Value{}, dilerr.NewGetError(
		fmt.Sprintf("There is no singletone with alias: %s", alias),
	)
}

func (di *dicon) GetTemporal(alias string, args ...interface{}) (interface{}, error) {
	container, err := di.getTemporal(alias, args...)
	if err != nil {
		return nil, err
	}
	return container.Interface(), nil
}

func (di *dicon) MustGetTemporal(alias string, args ...interface{}) interface{} {
	container, err := di.getTemporal(alias, args...)
	if err != nil {
		panic(err)
	}
	return container.Interface()
}

func (di *dicon) ProcessTemporal(alias string, container interface{}, args ...interface{}) error {
	c, err := di.getTemporal(alias, args...)
	v := reflect.ValueOf(container)
	if err == nil {
		v.Set(c)
	} else {
		v.Set(reflect.Zero(v.Type()))
	}

	return err
}

func (di *dicon) MustProcessTemporal(alias string, container interface{}, args ...interface{}) {
	c, err := di.getTemporal(alias, args...)
	v := reflect.ValueOf(container)
	if err != nil {
		panic(err)
	}
	v.Set(c)
}

// Get return services typed with some interface or construct and return service, if it is temporal.
func (di *dicon) getTemporal(alias string, args ...interface{}) (reflect.Value, error) {
	tempConstructor, ok := di.temporalByAlias[alias]
	if ok {
		argsIndex := 0
		creationResults, err := di.createService(tempConstructor, &argsIndex, args...)
		if err != nil {
			return reflect.Value{}, err
		}

		if len(creationResults) > 2 {
			return reflect.Value{}, dilerr.NewGetError(
				"temporal service creator returns more that 2 results",
			)
		}

		err, errIndex := checkHasError(creationResults)
		if errIndex != -1 && err != nil {
			return reflect.Value{}, err
		}

		return creationResults[0], nil
	}

	return reflect.Value{}, dilerr.NewGetError(
		fmt.Sprintf("There is no temporal service with alias: %s", alias),
	)
}
