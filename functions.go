package dilema

import (
	"dilema/dilerr"
	"reflect"
)

type callResults []reflect.Value

type CallResults interface {
	Process(values ...interface{}) error
	MustProcess(values ...interface{})
}

func (di *dicon) Run(function interface{}, args ...interface{}) (CallResults, error) {
	return di.processRunEvent(function, args...)
}

func (di *dicon) MustRun(function interface{}, args ...interface{}) CallResults {
	res, err := di.processRunEvent(function, args...)
	if err != nil {
		panic(err)
	}
	return res
}

func (di *dicon) processRunEvent(function interface{}, args ...interface{}) (CallResults, error) {
	operationIndex := <-di.operationIndexCh
	event := operationStartEvent{
		operationIndex: operationIndex,
		oType:          runOperation,
		event: runStartEvent{
			funcStartEvent: funcStartEvent{
				function: function,
				args:     args,
			},
		},
	}
	di.queueCh <- event

	for endEvent := range di.runEndCh {
		if endEvent.operationIndex == operationIndex {
			result := endEvent.result.(runEndEvent)
			return result.cr, result.err
		}
		di.runEndCh <- endEvent
	}
	return nil, dilerr.NewThreadError("run channel was closed")
}

func (di *dicon) Recover(function interface{}, args ...interface{}) (CallResults, error) {
	return di.processRecoverEvent(function, args...)
}

func (di *dicon) processRecoverEvent(function interface{}, args ...interface{}) (CallResults, error) {
	operationIndex := <-di.operationIndexCh
	event := operationStartEvent{
		operationIndex: operationIndex,
		oType:          recoverOperation,
		event: recoverStartEvent{
			funcStartEvent: funcStartEvent{
				function: function,
				args:     args,
			},
		},
	}
	di.queueCh <- event

	for endEvent := range di.recoverEndCh {
		if endEvent.operationIndex == operationIndex {
			result := endEvent.result.(recoverEndEvent)
			return result.cr, result.err
		}
		di.recoverEndCh <- endEvent
	}
	return nil, dilerr.NewThreadError("recover channel was closed")
}

func (di *dicon) recover(function interface{}, args ...interface{}) (cr CallResults, err error) {
	defer func() {
		if r := recover(); r != nil {
			cr, err = nil, r.(error)
		}
	}()

	cr, err = di.run(function, args...)
	if err != nil {
		panic(err)
	}

	return
}

func (di *dicon) RecoverAndClean(function interface{}, args ...interface{}) (CallResults, error) {
	return di.processRecoverAndCleanEvent(function, args...)
}

func (di *dicon) processRecoverAndCleanEvent(
	function interface{},
	args ...interface{},
) (CallResults, error) {
	operationIndex := <-di.operationIndexCh
	event := operationStartEvent{
		operationIndex: operationIndex,
		oType:          recoverAndCleanOperation,
		event: recoverAndCleanStartEvent{
			funcStartEvent: funcStartEvent{
				function: function,
				args:     args,
			},
		},
	}
	di.queueCh <- event

	for endEvent := range di.recoverAndCleanEndCh {
		if endEvent.operationIndex == operationIndex {
			result := endEvent.result.(recoverAndCleanEndEvent)
			return result.cr, result.err
		}
		di.recoverAndCleanEndCh <- endEvent
	}
	return nil, dilerr.NewThreadError("recover_and_clean channel was closed")
}

func (di *dicon) recoverAndClean(function interface{}, args ...interface{}) (cr CallResults, err error) {
	defer func() {
		if r := recover(); r != nil {
			di.clean()
			cr, err = nil, r.(error)
		}
	}()

	cr, err = di.run(function, args...)
	if err != nil {
		panic(err)
	}

	return
}

func (di *dicon) run(fun interface{}, args ...interface{}) (CallResults, error) {
	t, v := reflect.TypeOf(fun), reflect.ValueOf(fun)

	if t.Kind() != reflect.Func {
		return nil, dilerr.NewTypeError("unexpected fun type")
	}

	argsMap := make(map[reflect.Type][]reflect.Value)
	types := make([]reflect.Type, 0)
	for _, arg := range args {
		tArg, vArg := reflect.TypeOf(arg), reflect.ValueOf(arg)
		if arr, ok := argsMap[tArg]; ok {
			argsMap[tArg] = append(arr, vArg)
		} else {
			argsMap[tArg] = []reflect.Value{vArg}
		}
		types = append(types, tArg)
	}

	callArgs := make([]reflect.Value, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		tArg := t.In(i)
		if arr, ok := argsMap[tArg]; ok && len(arr) > 0 {
			callArgs[i] = arr[0]
			if len(arr) == 1 {
				delete(argsMap, tArg)
			} else {
				argsMap[tArg] = arr[1:]
			}
			continue
		}

		if tArg.Kind() == reflect.Interface {
			container, ok := di.getSingleToneByType(tArg)
			if ok {
				callArgs[i] = container
				continue
			}
			constuctor, ok := di.getTemporalByType(tArg)
			if ok {
				argsIndex := 0
				creationResults, err := di.createService(constuctor, &argsIndex, args...)
				if err != nil {
					return nil, err
				}
				errIndex, err := checkHasError(creationResults)
				if errIndex != -1 && err != nil {
					return nil, err
				}

				callArgs[i] = creationResults[0]
			}

			flag := false
			for _, tt := range types {
				if tt.Implements(tArg) {
					if arr, ok := argsMap[tt]; ok && len(arr) > 0 {
						callArgs[i] = arr[0]
						if len(arr) == 1 {
							delete(argsMap, tt)
						} else {
							argsMap[tt] = arr[1:]
						}

						flag = true
						break
					}
				}
			}
			if flag {
				continue
			}
		}

		if tArg.Kind() == reflect.Ptr &&
			tArg.Elem().Kind() == reflect.Struct {
			created, ok := di.createCorrectInStruct(tArg.Elem(), args...)
			if ok {
				callArgs[i] = created
				continue
			}
		}
		if tArg.Kind() == reflect.Struct {
			created, ok := di.createCorrectInStruct(tArg, args...)
			if ok {
				callArgs[i] = created.Elem()
				continue
			}
		}

		return nil, dilerr.NewTypeError("not enough arguments to call a function")
	}

	results := v.Call(callArgs)

	return callResults(results), nil
}

func (di *dicon) clean() {
	for _, destroyer := range di.getDestroyers() {
		destroyer.Call(nil)
	}
}

func (cr callResults) Process(values ...interface{}) error {
	return cr.process(values...)
}

func (cr callResults) MustProcess(values ...interface{}) {
	err := cr.process(values...)
	if err != nil {
		panic(err)
	}
}

func (cr callResults) process(values ...interface{}) error {
	crMap := make(map[reflect.Type][]reflect.Value)
	types := make([]reflect.Type, 0)
	for _, res := range cr {
		tRes := res.Type()
		if arr, ok := crMap[tRes]; ok {
			crMap[tRes] = append(arr, res)
		} else {
			crMap[tRes] = []reflect.Value{res}
		}
		types = append(types, tRes)
	}

	for _, val := range values {
		vVal := reflect.ValueOf(val)
		tVal := vVal.Type()
		elem := vVal.Elem()
		elemType := elem.Type()

		if val == nil {
			dilerr.NewTypeError("unexpected nil value")
		}
		if tVal.Kind() != reflect.Ptr {
			return dilerr.NewTypeError("expected ptr values")
		}
		if !elem.CanSet() {
			return dilerr.NewTypeError("agruments can't be setted")
		}
		if arr, ok := crMap[elemType]; ok {
			elem.Set(arr[0])
			if len(arr) == 1 {
				delete(crMap, elemType)
			} else {
				crMap[elemType] = arr[1:]
			}
			continue
		}

		flag := false
		for _, tt := range types {
			if elemType.Kind() == reflect.Interface && tt.Implements(elemType) {
				if arr, ok := crMap[tt]; ok && len(arr) > 0 {
					elem.Set(arr[0])
					if len(arr) == 1 {
						delete(crMap, elemType)
					} else {
						crMap[elemType] = arr[1:]
					}
					flag = true
					break
				}
			}
		}
		if flag {
			continue
		}
		return dilerr.NewTypeError("results and values are not comparable")
	}

	return nil
}
