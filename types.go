package dilema

import "reflect"

type operationType int

const (
	registerTemporaryOperation = iota
	registerSingleToneOperation
	registerFewOperation
	getTemporaryOperation
	getSingleToneOperation
	runOperation
	recoverOperation
	recoverAndCleanOperation
)

type operationStartEvent struct {
	operationCh chan operationEndEvent
	oType       int
	event       interface{}
}

type operationEndEvent struct {
	result interface{}
}

type registerTemporaryStartEvent struct {
	alias       string
	serviceInit interface{}
}

type registerSingleToneStartEvent struct {
	alias       string
	serviceInit interface{}
	args        []interface{}
}

type registerFewStartEvent struct {
	servicesInit map[string]interface{}
	args         []interface{}
}

type registerEndEvent struct {
	err error
}

type getTemporaryStartEvent struct {
	alias string
	args  []interface{}
}

type getSingleToneStartEvent struct {
	alias string
}

type getContainerEndEvent struct {
	container reflect.Value
	err       error
}

type funcStartEvent struct {
	function interface{}
	args     []interface{}
}

type funcEndEvent struct {
	cr  CallResults
	err error
}

type runStartEvent struct {
	funcStartEvent
}

type runEndEvent struct {
	funcEndEvent
}

type recoverStartEvent struct {
	funcStartEvent
}

type recoverEndEvent struct {
	funcEndEvent
}

type recoverAndCleanStartEvent struct {
	funcStartEvent
}

type recoverAndCleanEndEvent struct {
	funcEndEvent
}
