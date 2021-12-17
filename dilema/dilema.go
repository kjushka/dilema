package dilema

import "reflect"

type Dicon interface {
	ProvideTemp(serviceInit interface{}) error
	ProvideSingleTone(serviceInit interface{}, args ...interface{}) error
	ProvideAll(servicesInit ...interface{}) error
	Get(serviceAction interface{}, args ...interface{}) interface{}
	//GetWrapped(services... interface{}) interface{}
}

func Init() Dicon {
	return &dicon{
		temps:       make(map[reflect.Type]Constructor),
		singleTones: make(map[reflect.Type]reflect.Value),
	}
}
