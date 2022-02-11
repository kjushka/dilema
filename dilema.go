package dilema

import (
	"context"
)

type Dicon interface {
	RegisterTemporary(alias string, serviceInit interface{}) error
	MustRegisterTemporary(alias string, serviceInit interface{})
	RegisterSingletone(alias string, serviceInit interface{}, args ...interface{}) error
	MustRegisterSingletone(alias string, serviceInit interface{}, args ...interface{})
	RegisterFew(servicesInit map[string]interface{}, args ...interface{}) error
	MustRegisterFew(servicesInit map[string]interface{}, args ...interface{})
	GetSingletone(alias string) (interface{}, error)
	MustGetSingletone(alias string) interface{}
	ProcessSingletone(alias string, container interface{}) error
	MustProcessSingletone(alias string, container interface{})
	GetTemporary(alias string, args ...interface{}) (interface{}, error)
	MustGetTemporary(alias string, args ...interface{}) interface{}
	ProcessTemporary(alias string, container interface{}, args ...interface{}) error
	MustProcessTemporary(alias string, container interface{}, args ...interface{})
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
		temporaryStore: newTemporaryStore(),
		singleToneStore: newSingleToneStore(),
		destroyerStore: newDestroyerStore(),

		queueStore: newQueueStore(),

		operationStartCh: make(chan operationStartEvent),
		queueCh: make(chan operationStartEvent),
		exitCh: make(chan struct{}),

		ctx: context.Background(),
	}

	go di.goQueueWriter()
	go di.goQueueReader()

	return di
}
