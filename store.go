package dilema

import (
	"container/list"
	"reflect"
	"sync"
)

type queueStore struct {
	sync.Mutex
	queue *list.List
}

func newQueueStore() *queueStore {
	return &queueStore{
		queue: list.New(),
	}
}

func (qs *queueStore) pushEventBack(event operationStartEvent) {
	qs.Lock()
	defer qs.Unlock()

	qs.queue.PushBack(event)
}

func (qs *queueStore) popEvent() operationStartEvent {
	qs.Lock()
	defer qs.Unlock()

	el := qs.queue.Front()
	return qs.queue.Remove(el).(operationStartEvent)
}

func (qs *queueStore) queueLen() int {
	qs.Lock()
	defer qs.Unlock()

	return qs.queue.Len()
}

type temporalStore struct {
	sync.RWMutex
	temporalByAlias map[string]reflect.Value
	temporalByType  map[reflect.Type]reflect.Value
}

func newTemporalStore() *temporalStore {
	return &temporalStore{
		temporalByAlias: make(map[string]reflect.Value),
		temporalByType:  make(map[reflect.Type]reflect.Value),
	}
}

func (ts *temporalStore) addTemporal(alias string, v reflect.Value, t reflect.Type) {
	ts.Lock()
	defer ts.Unlock()

	ts.temporalByAlias[alias] = v
	ts.temporalByType[t] = v
}

func (ts *temporalStore) getTemporalByAlias(alias string) (reflect.Value, bool) {
	ts.RLock()
	defer ts.RUnlock()

	temporal, ok := ts.temporalByAlias[alias]
	return temporal, ok
}

func (ts *temporalStore) getTemporalByType(t reflect.Type) (reflect.Value, bool) {
	ts.RLock()
	defer ts.RUnlock()

	temporal, ok := ts.temporalByType[t]
	return temporal, ok
}

type singleToneStore struct {
	sync.RWMutex
	singleTonesByAlias map[string]reflect.Value
	singleTonesByType  map[reflect.Type]reflect.Value
}

func newSingleToneStore() *singleToneStore {
	return &singleToneStore{
		singleTonesByAlias: make(map[string]reflect.Value),
		singleTonesByType:  make(map[reflect.Type]reflect.Value),
	}
}

func (ss *singleToneStore) addSingleTone(alias string, v reflect.Value, t reflect.Type) {
	ss.Lock()
	defer ss.Unlock()

	ss.singleTonesByAlias[alias] = v
	ss.singleTonesByType[t] = v
}

func (ss *singleToneStore) getSingleToneByAlias(alias string) (reflect.Value, bool) {
	ss.RLock()
	defer ss.RUnlock()

	temporal, ok := ss.singleTonesByAlias[alias]
	return temporal, ok
}

func (ss *singleToneStore) getSingleToneByType(t reflect.Type) (reflect.Value, bool) {
	ss.RLock()
	defer ss.RUnlock()

	temporal, ok := ss.singleTonesByType[t]
	return temporal, ok
}

type destroyerStore struct {
	sync.Mutex
	destroyers []reflect.Value
}

func newDestroyerStore() *destroyerStore {
	return &destroyerStore{
		destroyers: make([]reflect.Value, 0),
	}
}

func (ds *destroyerStore) addDestroyer(destroyer reflect.Value) {
	ds.Lock()
	defer ds.Unlock()

	ds.destroyers = append(ds.destroyers, destroyer)
}

func (ds *destroyerStore) getDestroyers() []reflect.Value {
	return ds.destroyers
}
