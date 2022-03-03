package dilema

import (
	"context"
	"dilema/dilerr"
	"fmt"
	"reflect"
)

type dicon struct {
	*temporaryStore
	*singleToneStore
	*destroyablesStore

	ctx context.Context
}

func (di *dicon) RegisterTemporary(alias string, serviceInit interface{}) error {
	return di.registerTemporary(alias, serviceInit)
}

func (di *dicon) MustRegisterTemporary(alias string, serviceInit interface{}) {
	err := di.registerTemporary(alias, serviceInit)
	if err != nil {
		panic(err)
	}
}

// registerTemporary provides new service, which will be initialized when
// you call Get method and be destroyed with GC after work will be done
func (di *dicon) registerTemporary(alias string, serviceInit interface{}) error {
	if _, ok := di.getTemporaryByAlias(alias); ok {
		return dilerr.GetAlreadyExistError(alias)
	}
	_, v, err := checkProvidedTypeIsCreator(serviceInit)
	if err != nil {
		return err
	}

	di.addTemporary(alias, v)
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
	if _, ok := di.getSingleToneByAlias(alias); ok {
		return dilerr.GetAlreadyExistError(alias)
	}
	t, v, err := checkProvidedTypeIsCreator(serviceInit)
	if err != nil {
		return err
	}

	creationResults, err := di.createService(v, args...)
	if err != nil {
		return err
	}

	if t.NumOut() == 2 && len(creationResults) == 2 {
		err, ok := checkIsError(creationResults[1])
		if ok {
			return err
		} 
	}

	di.addSingleTone(alias, creationResults[0], t)
	if destroyable, ok := creationResults[0].Interface().(Destroyable); ok {
		di.addDestroyable(destroyable)
	}

	return nil
}

// createService creates instance of service, which interface return from provided func
func (di *dicon) createService(
	v reflect.Value,
	args ...interface{},
) (callResults, error) {
	return di.run(v, args...)
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
	if err != nil {
		return err
	}
	err = processValue(c, container)

	return err
}

func (di *dicon) MustProcessSingletone(alias string, container interface{}) {
	c, err := di.getSingletone(alias)
	if err != nil {
		panic(err)
	}
	err = processValue(c, container)
	if err != nil {
		panic(err)
	}
}

func (di *dicon) getSingletone(alias string) (reflect.Value, error) {
	singleTone, ok := di.getSingleToneByAlias(alias)
	if ok {
		return singleTone, nil
	}
	return reflect.Value{}, dilerr.NewGetError(
		fmt.Sprintf("there is no singletone with alias: %s", alias),
	)
}

func (di *dicon) GetTemporary(alias string, args ...interface{}) (interface{}, error) {
	container, err := di.getTemporary(alias, args...)
	if err != nil {
		return nil, err
	}
	return container.Interface(), nil
}

func (di *dicon) MustGetTemporary(alias string, args ...interface{}) interface{} {
	container, err := di.getTemporary(alias, args...)
	if err != nil {
		panic(err)
	}
	return container.Interface()
}

func (di *dicon) ProcessTemporary(alias string, container interface{}, args ...interface{}) error {
	c, err := di.getTemporary(alias, args...)
	if err != nil {
		return err
	}
	err = processValue(c, container)

	return err
}

func (di *dicon) MustProcessTemporary(alias string, container interface{}, args ...interface{}) {
	c, err := di.getTemporary(alias, args...)
	if err != nil {
		panic(err)
	}
	err = processValue(c, container)
	if err != nil {
		panic(err)
	}
}

// Get return services typed with some interface or construct and return service, if it is temporary.
func (di *dicon) getTemporary(alias string, args ...interface{}) (reflect.Value, error) {
	tempConstructor, ok := di.getTemporaryByAlias(alias)
	if ok {
		creationResults, err := di.createService(tempConstructor, args...)
		if err != nil {
			return reflect.Value{}, err
		}

		if tempConstructor.Type().NumOut() == 2 && len(creationResults) == 2 {
			err, ok := checkIsError(creationResults[1])
			if ok {
				return reflect.Value{}, err
			} 
		}

		return creationResults[0], nil
	}

	return reflect.Value{}, dilerr.NewGetError(
		fmt.Sprintf("there is no temporary service with alias: %s", alias),
	)
}
