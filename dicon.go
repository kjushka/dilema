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
	*destroyerStore

	*queueStore

	operationStartCh chan operationStartEvent
	queueCh          chan operationStartEvent
	exitCh           chan struct{}

	ctx context.Context
}

func (di *dicon) RegisterTemporary(alias string, serviceInit interface{}) error {
	return di.processRegisterTemporaryEvent(alias, serviceInit)
}

func (di *dicon) MustRegisterTemporary(alias string, serviceInit interface{}) {
	err := di.processRegisterTemporaryEvent(alias, serviceInit)
	if err != nil {
		panic(err)
	}
}

func (di *dicon) processRegisterTemporaryEvent(alias string, serviceInit interface{}) error {
	operationCh := make(chan operationEndEvent)
	startEvent := operationStartEvent{
		operationCh: operationCh,
		oType:       registerTemporaryOperation,
		event: registerTemporaryStartEvent{
			alias:       alias,
			serviceInit: serviceInit,
		},
	}
	di.queueCh <- startEvent

	endEvent := <-operationCh
	close(operationCh)

	return endEvent.result.(registerEndEvent).err
}

// registerTemporary provides new service, which will be initialized when
// you call Get method and be destroyed with GC after work will be done
func (di *dicon) registerTemporary(alias string, serviceInit interface{}) error {
	if _, ok := di.getTemporaryByAlias(alias); ok {
		return dilerr.GetAlreadyExistError(alias)
	}
	t, v, err := checkProvidedTypeIsCreator(serviceInit)
	if err != nil {
		return err
	}

	di.addTemporary(alias, v, t)
	return nil
}

func (di *dicon) RegisterSingletone(
	alias string,
	serviceInit interface{},
	args ...interface{},
) error {
	return di.processRegisterSingleToneEvent(alias, serviceInit, args...)
}

func (di *dicon) MustRegisterSingletone(
	alias string,
	serviceInit interface{},
	args ...interface{},
) {
	err := di.processRegisterSingleToneEvent(alias, serviceInit, args...)
	if err != nil {
		panic(err)
	}
}

func (di *dicon) processRegisterSingleToneEvent(
	alias string,
	serviceInit interface{},
	args ...interface{},
) error {
	operationCh := make(chan operationEndEvent)
	event := operationStartEvent{
		operationCh: operationCh,
		oType:       registerSingleToneOperation,
		event: registerSingleToneStartEvent{
			alias:       alias,
			serviceInit: serviceInit,
			args:        args,
		},
	}
	di.queueCh <- event

	endEvent := <-operationCh
	close(operationCh)

	return endEvent.result.(registerEndEvent).err
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

	argsIndex := 0
	creationResults, err := di.createService(v, &argsIndex, args...)
	if err != nil {
		return err
	}
	destroyerIndex, err := di.checkCreationResults(creationResults)
	if err != nil {
		return err
	}

	di.addSingleTone(alias, v, t)
	if destroyerIndex != -1 {
		di.addDestroyer(creationResults[destroyerIndex])
	}

	return nil
}

func (di *dicon) RegisterFew(servicesInit map[string]interface{}, args ...interface{}) error {
	return di.processRegisterFewEvent(servicesInit, args...)
}

func (di *dicon) MustRegisterFew(servicesInit map[string]interface{}, args ...interface{}) {
	err := di.processRegisterFewEvent(servicesInit, args...)
	if err != nil {
		panic(err)
	}
}

func (di *dicon) processRegisterFewEvent(servicesInit map[string]interface{}, args ...interface{}) error {
	operationCh := make(chan operationEndEvent)
	event := operationStartEvent{
		operationCh: operationCh,
		oType:       registerFewOperation,
		event: registerFewStartEvent{
			servicesInit: servicesInit,
			args:         args,
		},
	}
	di.queueCh <- event

	endEvent := <-operationCh
	close(operationCh)

	return endEvent.result.(registerEndEvent).err
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

		var destroyer interface{}
		if destroyerIndex != -1 {
			destroyer = creationResults[destroyerIndex]
		} else {
			destroyer = nil
		}

		services = append(services, ta{
			a,
			creationResults[0].Type(),
			creationResults[0],
			destroyer,
		})
	}

	for _, service := range services {
		di.addSingleTone(service.a, service.v, service.t)
		if service.d != nil {
			destroyer := service.d.(reflect.Value)
			di.addDestroyer(destroyer)
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
	temp, ok := di.getTemporaryByType(paramT)
	if ok {
		creationResults, err := di.createService(temp, argsIndex, args...)
		if err != nil {
			return reflect.Value{}, false, err
		}
		return creationResults[0], true, nil
	}
	singleTone, ok := di.getSingleToneByType(paramT)
	if ok {
		return singleTone, ok, nil
	}
	return reflect.Value{}, ok, nil
}

func (di *dicon) GetSingletone(alias string) (interface{}, error) {
	container, err := di.processGetSingleToneEvent(alias)
	if err != nil {
		return nil, err
	}
	return container.Interface(), nil
}

func (di *dicon) MustGetSingletone(alias string) interface{} {
	container, err := di.processGetSingleToneEvent(alias)
	if err != nil {
		panic(err)
	}
	return container.Interface()
}

func (di *dicon) ProcessSingletone(alias string, container interface{}) error {
	c, err := di.processGetSingleToneEvent(alias)
	if err != nil {
		return err
	}
	err = processValue(c, container)

	return err
}

func (di *dicon) MustProcessSingletone(alias string, container interface{}) {
	c, err := di.processGetSingleToneEvent(alias)
	if err != nil {
		panic(err)
	}
	err = processValue(c, container)
	if err != nil {
		panic(err)
	}
}

func (di *dicon) processGetSingleToneEvent(alias string) (reflect.Value, error) {
	operationCh := make(chan operationEndEvent)
	event := operationStartEvent{
		operationCh: operationCh,
		oType:       getSingleToneOperation,
		event: getSingleToneStartEvent{
			alias: alias,
		},
	}
	di.queueCh <- event

	endEvent := <-operationCh
	close(operationCh)
	result := endEvent.result.(getContainerEndEvent)

	return result.container, result.err
}

func (di *dicon) getSingletone(alias string) (reflect.Value, error) {
	singleTone, ok := di.getSingleToneByAlias(alias)
	if ok {
		return singleTone, nil
	}
	return reflect.Value{}, dilerr.NewGetError(
		fmt.Sprintf("There is no singletone with alias: %s", alias),
	)
}

func (di *dicon) GetTemporary(alias string, args ...interface{}) (interface{}, error) {
	container, err := di.processGetTemporaryEvent(alias, args...)
	if err != nil {
		return nil, err
	}
	return container.Interface(), nil
}

func (di *dicon) MustGetTemporary(alias string, args ...interface{}) interface{} {
	container, err := di.processGetTemporaryEvent(alias, args...)
	if err != nil {
		panic(err)
	}
	return container.Interface()
}

func (di *dicon) ProcessTemporary(alias string, container interface{}, args ...interface{}) error {
	c, err := di.processGetTemporaryEvent(alias, args...)
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

func (di *dicon) processGetTemporaryEvent(alias string, args ...interface{}) (reflect.Value, error) {
	operationCh := make(chan operationEndEvent)
	event := operationStartEvent{
		operationCh: operationCh,
		oType:       getTemporaryOperation,
		event: getTemporaryStartEvent{
			alias: alias,
			args:  args,
		},
	}
	di.queueCh <- event

	endEvent := <-operationCh
	close(operationCh)
	result := endEvent.result.(getContainerEndEvent)

	return result.container, result.err
}

// Get return services typed with some interface or construct and return service, if it is temporary.
func (di *dicon) getTemporary(alias string, args ...interface{}) (reflect.Value, error) {
	tempConstructor, ok := di.getTemporaryByAlias(alias)
	if ok {
		argsIndex := 0
		creationResults, err := di.createService(tempConstructor, &argsIndex, args...)
		if err != nil {
			return reflect.Value{}, err
		}

		if len(creationResults) > 2 {
			return reflect.Value{}, dilerr.NewGetError(
				"temporary service creator returns more that 2 results",
			)
		}

		errIndex, err := checkHasError(creationResults)
		if errIndex != -1 && err != nil {
			return reflect.Value{}, err
		}

		return creationResults[0], nil
	}

	return reflect.Value{}, dilerr.NewGetError(
		fmt.Sprintf("There is no temporary service with alias: %s", alias),
	)
}
