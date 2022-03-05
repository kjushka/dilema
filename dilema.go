package dilema

import (
	"context"
)

type Dicon interface {
	// RegisterTemporary registers a temporary container where serviceInit is a function-constructor
	// that returns either an interface value or an interface value and an error. 
	// The method returns a type mismatch error, if one was found. 
	// If you try to register a new container with a previously used alias, an error will also be returned
	RegisterTemporary(alias string, serviceInit interface{}) error
	// MustRegisterTemporary registers a temporary container where serviceInit is a function-constructor
	// that returns either an interface value or an interface value and an error. 
	// The method panics with mismatch error, if one was found. 
	// If you try to register a new container with a previously used alias, method also will panic.
	// If the container creation fails, the panic(err) method will be called
	MustRegisterTemporary(alias string, serviceInit interface{})
	// RegisterSingletone registers a permanent container, where serviceInit is a function 
	// that returns either an interface value or an interface value and an error, 
	// agrs are the arguments necessary to create a service, 
	// passed in the order in which the constructor function accepts them. 
	// The method returns a type mismatch error, if one was found, 
	// and errors that occur when creating a container. 
	// When trying to register a new container with a previously used alias, 
	// an error will also be returned.
	RegisterSingletone(alias string, serviceInit interface{}, args ...interface{}) error
	// MustRegisterSingletone registers a permanent container, where serviceInit is a function 
	// that returns either an interface value or an interface value and an error, 
	// agrs are the arguments necessary to create a service, 
	// passed in the order in which the constructor function accepts them. 
	// The method panics when types are mismatch or if errors occur when creating a container. 
	// When trying to register a new container with a previously used alias, it also will panic.
	MustRegisterSingletone(alias string, serviceInit interface{}, args ...interface{})
	// GetSingletone returns a previously registered permanent container by the passed alias. 
	// If the container is not found, an error will be returned.
	GetSingletone(alias string) (interface{}, error)
	// MustGetSingletone returns a previously registered permanent container by the passed alias. 
	// If the container is not found, method will panic.
	MustGetSingletone(alias string) interface{}
	// ProcessSingleTone allows you to substitute a permanent container registered earlier 
	// by the passed alias into a variable passed by reference as an argument 'container'. 
	// If the container is not found or substitution is not possible, an error will be returned.
	ProcessSingletone(alias string, container interface{}) error
	// MustProcessSingleTone allows you to substitute a permanent container registered earlier 
	// by the passed alias into a variable passed by reference as an argument 'container'. 
	// If the container is not found or substitution is not possible, method will panic.
	MustProcessSingletone(alias string, container interface{})
	// GetTemporary returns a previously registered temporary container by the passed alias. 
	// To create a container, you must pass the arguments in the order 
	// in which the constructor function accepts them. 
	// If the container is not found or the creation of the container failed, an error will be returned.
	GetTemporary(alias string, args ...interface{}) (interface{}, error)
	// MustGetTemporary returns a previously registered temporary container by the passed alias. 
	// To create a container, you must pass the arguments in the order 
	// in which the constructor function accepts them. 
	// If the container is not found or the creation of the container failed, method will panic.
	MustGetTemporary(alias string, args ...interface{}) interface{}
	// ProcessTemporary allows you to substitute a permanent container 
	// registered earlier by the passed alias into a variable passed by reference as an argument 'container'. 
	// To create a container, you must pass the arguments in the order 
	// in which the constructor function accepts them. 
	// If the container is not found, substitution is not possible, or container creation failed, an error will be returned.
	ProcessTemporary(alias string, container interface{}, args ...interface{}) error
	// MustProcessTemporary allows you to substitute a permanent container 
	// registered earlier by the passed alias into a variable passed by reference as an argument 'container'. 
	// To create a container, you must pass the arguments in the order 
	// in which the constructor function accepts them. 
	// If the container is not found, substitution is not possible, or container creation failed, method will panic.
	MustProcessTemporary(alias string, container interface{}, args ...interface{})
	// ProcessStruct "collects" the struct passed to the method as an argument, 
	// provided that the fields are public and have the types of previously registered permanent containers. 
	// Also, the fields of the structure can have the 'dilema:"container name"' tag specified for reliable substitution.
	// If the structure assembly or substitution failed, an error will be returned.
	ProcessStruct(str interface{}) error
	// MustProcessStruct "collects" the struct passed to the method as an argument, 
	// provided that the fields are public and have the types of previously registered permanent containers. 
	// Also, the fields of the structure can have the 'dilema:"container name"' tag specified for reliable substitution.
	// If the structure assembly or substitution failed, method will panic.
	MustProcessStruct(str interface{})
	// Run calls the function passed as the first argument. 
	// The arguments of the function must be passed in the order 
	// in which they are required to call the function. 
	// Returns CallResults and an error, if any occures when the function was started.
	Run(function interface{}, args ...interface{}) (CallResults, error)
	// Run calls the function passed as the first argument. 
	// The arguments of the function must be passed in the order 
	// in which they are required to call the function. 
	// Returns CallResults. If any error occures when the function was started, method will panic.
	MustRun(function interface{}, args ...interface{}) CallResults
	// Recover calls the function passed as the first argument. 
	// The arguments of the function must be passed in the order 
	// in which they are required to call the function. 
	// Returns Call Results and an error, if any occures when the function was started. 
	// In case of panic inside the function, this method processes and returns the error that occurred.
	Recover(function interface{}, args ...interface{}) (cr CallResults, err error)
	// RecoverAndClean calls the function passed as the first argument. 
	// The arguments of the function must be passed in the order 
	// in which they are required to call the function. 
	// Returns Call Results and an error, if any occures when the function was started. 
	// In case of panic inside the function, this method processes and returns the error 
	// that has occurred, and also calls the Destroy() method for registered permanent containers.
	RecoverAndClean(function interface{}, args ...interface{}) (cr CallResults, err error)
	// Ctx allows you to get the context of the DI container.
	Ctx() context.Context
	// SetCtx allows you to set a new context in the DI-container.
	SetCtx(ctx context.Context)
	// AddToCtx allows you to add a new value to the context with the specified alias.
	AddToCtx(alias string, value interface{})
	// GetFromCtx allows you to get the value from the context by the passed alias.
	GetFromCtx(alias string) interface{}
}

func Init() Dicon {
	di := &dicon{
		temporaryStore:    newTemporaryStore(),
		singleToneStore:   newSingleToneStore(),
		destroyablesStore: newDestroyablesStore(),

		ctx: context.Background(),
	}

	return di
}
