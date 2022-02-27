package dilema

import (
	"reflect"
	"sync"
)

type temporaryStore struct {
	sync.RWMutex
	temporaryByAlias map[string]reflect.Value
}

func newTemporaryStore() *temporaryStore {
	return &temporaryStore{
		temporaryByAlias: make(map[string]reflect.Value),
	}
}

func (ts *temporaryStore) addTemporary(alias string, v reflect.Value) {
	ts.Lock()
	defer ts.Unlock()

	ts.temporaryByAlias[alias] = v
}

func (ts *temporaryStore) getTemporaryByAlias(alias string) (reflect.Value, bool) {
	ts.RLock()
	defer ts.RUnlock()

	temporary, ok := ts.temporaryByAlias[alias]
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

type destroyablesStore struct {
	sync.Mutex
	destroyables []Destroyable
}

func newDestroyablesStore() *destroyablesStore {
	return &destroyablesStore{
		destroyables: make([]Destroyable, 0),
	}
}

func (ds *destroyablesStore) addDestroyable(destroyable Destroyable) {
	ds.Lock()
	defer ds.Unlock()

	ds.destroyables = append(ds.destroyables, destroyable)
}

func (ds *destroyablesStore) getDestroyables() []Destroyable {
	return ds.destroyables
}
