package dilema

import (
	"reflect"
	"sync"
)

type Dicon interface {
	ProvideTemporal(serviceInit interface{}) error
	ProvideSingleTone(serviceInit interface{}, args ...interface{}) error
	ProvideAll(servicesInit ...interface{}) error
	Get(serviceAction interface{}, args ...interface{}) interface{}
	//GetFromUnion(union interface{}, args ...interface{}) interface{}
}

func Init() Dicon {
	return &dicon{
		temps:       make(map[reflect.Type]Constructor),
		singleTones: make(map[reflect.Type]reflect.Value),
		functions:   make(map[reflect.Method]actionToType),
		cache:       make(map[reflect.Type]reflect.Value),

		mutex: &sync.Mutex{},
	}
}
