package dilema

import (
	"github.com/kjushka/dilema/dilerr"
	"reflect"
)

type callResults []reflect.Value

type CallResults interface {
	Process(values ...interface{}) error
	MustProcess(values ...interface{})
}

func (di *dicon) Run(function interface{}, args ...interface{}) (CallResults, error) {
	v := reflect.ValueOf(function)
	return di.run(v, args...)
}

func (di *dicon) MustRun(function interface{}, args ...interface{}) CallResults {
	v := reflect.ValueOf(function)
	res, err := di.run(v, args...)
	if err != nil {
		panic(err)
	}
	return res
}

func (di *dicon) Recover(function interface{}, args ...interface{}) (cr CallResults, err error) {
	defer func() {
		if r := recover(); r != nil {
			cr, err = nil, r.(error)
		}
	}()

	v := reflect.ValueOf(function)
	cr, err = di.run(v, args...)
	if err != nil {
		panic(err)
	}

	return
}

func (di *dicon) RecoverAndClean(function interface{}, args ...interface{}) (cr CallResults, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = di.clean()
			if err != nil {
				cr = nil
				return
			}
			cr, err = nil, r.(error)
		}
	}()

	v := reflect.ValueOf(function)
	cr, err = di.run(v, args...)
	if err != nil {
		panic(err)
	}

	return
}

func (di *dicon) run(fun reflect.Value, args ...interface{}) (cr callResults, err error) {
	t := fun.Type()

	if t.Kind() != reflect.Func {
		return nil, dilerr.NewTypeError("unexpected fun type")
	}

	callArgs := make([]reflect.Value, 0)
	if len(args) == 0 {
		results := fun.Call(callArgs)
		return callResults(results), nil
	}

	argIndex := 0
	currentArgument := reflect.ValueOf(args[argIndex])
	updateArgument := func() {
		argIndex++
		if argIndex != len(args) {
			currentArgument = reflect.ValueOf(args[argIndex])
		}
	}

	for i := 0; i < t.NumIn(); i++ {
		tIn := t.In(i)
		if tIn == currentArgument.Type() {
			callArgs = append(callArgs, currentArgument)
			updateArgument()
			continue
		}

		if tIn.Kind() == reflect.Interface {
			if currentArgument.Type().Implements(tIn) {
				callArgs = append(callArgs, currentArgument)
				updateArgument()
				continue
			}

			container, ok := di.getSingleToneByType(tIn)
			if ok {
				callArgs = append(callArgs, container)
				continue
			}
		}

		if tIn.Kind() == reflect.Ptr &&
			tIn.Elem().Kind() == reflect.Struct {
			created, ok := di.createInStruct(tIn.Elem())
			if ok {
				callArgs = append(callArgs, created)
				continue
			}
		}
		if tIn.Kind() == reflect.Struct {
			created, ok := di.createInStruct(tIn)
			if ok {
				callArgs = append(callArgs, created)
				continue
			}
		}

		return nil, dilerr.NewTypeError("not enough arguments to call a function")
	}

	results := fun.Call(callArgs)

	return callResults(results), nil
}

func (di *dicon) clean() error {
	for _, destroyable := range di.getDestroyables() {
		err := destroyable.Destroy()
		if err != nil {
			return err
		}
	}
	return nil
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

func (cr callResults) process(values ...interface{}) (err error) {
	if len(cr) != len(values) {
		return dilerr.NewProcessError("expected same count of values as results length")
	}

	for i := 0; i < len(cr); i++ {
		err = processValue(cr[i], values[i])
		if err != nil {
			return err
		}
	}

	err = nil
	return
}
