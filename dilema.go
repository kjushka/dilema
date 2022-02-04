package dilema

import (
	"reflect"
	"sync"
)

type Dicon interface {
	RegisterTemporal(alias string, serviceInit interface{}) error
	MustRegisterTemporal(alias string, serviceInit interface{})
	RegisterSingletone(alias string, serviceInit interface{}, args ...interface{}) error
	MustRegisterSingletone(alias string, serviceInit interface{}, args ...interface{})
	RegisterFew(servicesInit map[string]interface{}, args ...interface{}) error
	MustRegisterFew(servicesInit map[string]interface{}, args ...interface{})
	GetSingletone(alias string) (interface{}, error)
	MustGetSingletone(alias string) interface{}
	ProcessSingletone(alias string, container interface{}) error
	MustProcessSingletone(alias string, container interface{})
	GetTemporal(alias string, args ...interface{}) (interface{}, error)
	MustGetTemporal(alias string, args ...interface{}) interface{}
	ProcessTemporal(alias string, container interface{}, args ...interface{}) error
	MustProcessTemporal(alias string, container interface{}, args ...interface{})
}

func Init() Dicon {
	return &dicon{
		temporalByAlias:    make(map[string]reflect.Value),
		temporalByType:     make(map[reflect.Type]reflect.Value),
		singletonesByAlias: make(map[string]reflect.Value),
		singletonesByType:  make(map[reflect.Type]reflect.Value),

		destroyers: make([]reflect.Value, 0),
		cache:      make(map[reflect.Type]reflect.Value),

		mutex: &sync.Mutex{},
	}
}
