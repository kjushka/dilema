package dilema

import (
	"context"
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
	Run(function interface{}, args ...interface{}) (CallResults, error)
	MustRun(function interface{}, args ...interface{}) CallResults
	Recover(function interface{}, args ...interface{}) (cr CallResults, err error)
	RecoverAndClean(function interface{}, args ...interface{}) (cr CallResults, err error)
	Ctx() context.Context
	SetCtx(ctx context.Context)
	AddToCtx(alias string, value interface{})
	GetFromCtx(alias string) interface{}
}

func Init() Dicon {
	di := &dicon{
		temporalStore: newTemporalStore(),
		singleToneStore: newSingleToneStore(),
		destroyerStore: newDestroyerStore(),

		queueStore: newQueueStore(),

		operationIndexCh: make(chan uint64),

		operationStartCh: make(chan operationStartEvent),
		operationEndChansStore: newOperationEndChansStore(),
		queueCh: make(chan operationStartEvent),
		exitCh: make(chan struct{}),

		ctx: context.Background(),
	}

	go di.goOperationIndexProvider()
	go di.goQueueWriter()
	go di.goQueueReader()

	return di
}
