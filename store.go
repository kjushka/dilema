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

type temporaryStore struct {
	sync.RWMutex
	temporaryByAlias map[string]reflect.Value
	temporaryByType  map[reflect.Type]reflect.Value
}

func newTemporaryStore() *temporaryStore {
	return &temporaryStore{
		temporaryByAlias: make(map[string]reflect.Value),
		temporaryByType:  make(map[reflect.Type]reflect.Value),
	}
}

func (ts *temporaryStore) addTemporary(alias string, v reflect.Value, t reflect.Type) {
	ts.Lock()
	defer ts.Unlock()

	ts.temporaryByAlias[alias] = v
	ts.temporaryByType[t] = v
}

func (ts *temporaryStore) getTemporaryByAlias(alias string) (reflect.Value, bool) {
	ts.RLock()
	defer ts.RUnlock()

	temporary, ok := ts.temporaryByAlias[alias]
	return temporary, ok
}

func (ts *temporaryStore) getTemporaryByType(t reflect.Type) (reflect.Value, bool) {
	ts.RLock()
	defer ts.RUnlock()

	temporary, ok := ts.temporaryByType[t]
	return temporary, ok
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

	singletone, ok := ss.singleTonesByAlias[alias]
	return singletone, ok
}

func (ss *singleToneStore) getSingleToneByType(t reflect.Type) (reflect.Value, bool) {
	ss.RLock()
	defer ss.RUnlock()

	singletone, ok := ss.singleTonesByType[t]
	return singletone, ok
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
