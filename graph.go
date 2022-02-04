package dilema

import (
	"reflect"
	"sync"
)

type graph struct {
	nodes []*node
	edges map[*node][]*node
	mutex sync.RWMutex
}

type node struct {
	cType containerType
	vType reflect.Type
	alias string
}

func newGraph() *graph {
	return &graph{}
}

func (g *graph) addNode(n *node) {
	//checkNoCycle
	g.mutex.Lock()
	g.nodes = append(g.nodes, n)
	g.mutex.Unlock()
}

func (g *graph) addEdge(n *node) {
	g.nodes = append(g.nodes, n);
}
