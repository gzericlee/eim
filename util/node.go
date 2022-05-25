package util

import (
	"crypto/sha1"
	"math"
	"sort"
	"strconv"
	"sync"
)

const (
	DefaultVirtualSpots = 400
)

type node struct {
	nodeKey   string
	spotValue uint32
}

type instance struct {
	weight int
	object interface{}
}

type nodesArray []node

func (its nodesArray) Len() int           { return len(its) }
func (its nodesArray) Less(i, j int) bool { return its[i].spotValue < its[j].spotValue }
func (its nodesArray) Swap(i, j int)      { its[i], its[j] = its[j], its[i] }
func (its nodesArray) Sort()              { sort.Sort(its) }

type NodeManager struct {
	virtualSpots int
	nodes        nodesArray
	instances    map[string]*instance
	mu           sync.RWMutex
}

func NewNodeManager(spots int) *NodeManager {
	if spots == 0 {
		spots = DefaultVirtualSpots
	}

	h := &NodeManager{
		virtualSpots: spots,
		instances:    make(map[string]*instance),
	}
	return h
}

func (its *NodeManager) AddNode(nodeKey string, weight int, object interface{}) {
	its.mu.Lock()
	defer its.mu.Unlock()
	its.instances[nodeKey] = &instance{
		weight: weight,
		object: object,
	}
	its.generate()
}

func (its *NodeManager) RemoveNode(nodeKey string) {
	its.mu.Lock()
	defer its.mu.Unlock()
	delete(its.instances, nodeKey)
	its.generate()
}

func (its *NodeManager) GetNodes() []string {
	var keys []string
	for key, _ := range its.instances {
		keys = append(keys, key)
	}
	return keys
}

func (its *NodeManager) generate() {
	var totalW int
	for _, ins := range its.instances {
		totalW += ins.weight
	}

	totalVirtualSpots := its.virtualSpots * len(its.instances)
	its.nodes = nodesArray{}

	for nodeKey, ins := range its.instances {
		spots := int(math.Floor(float64(ins.weight) / float64(totalW) * float64(totalVirtualSpots)))
		for i := 1; i <= spots; i++ {
			hash := sha1.New()
			hash.Write([]byte(nodeKey + ":" + strconv.Itoa(i)))
			hashBytes := hash.Sum(nil)
			n := node{
				nodeKey:   nodeKey,
				spotValue: genValue(hashBytes[6:10]),
			}
			its.nodes = append(its.nodes, n)
			hash.Reset()
		}
	}
	its.nodes.Sort()
}

func genValue(bs []byte) uint32 {
	if len(bs) < 4 {
		return 0
	}
	v := (uint32(bs[3]) << 24) | (uint32(bs[2]) << 16) | (uint32(bs[1]) << 8) | (uint32(bs[0]))
	return v
}

func (its *NodeManager) GetNode(val string) (string, interface{}) {
	if len(its.nodes) == 0 {
		return "", nil
	}

	hash := sha1.New()
	hash.Write([]byte(val))
	hashBytes := hash.Sum(nil)
	v := genValue(hashBytes[6:10])
	i := sort.Search(len(its.nodes), func(i int) bool { return its.nodes[i].spotValue >= v })

	if i == len(its.nodes) {
		i = 0
	}
	key := its.nodes[i].nodeKey

	return key, its.instances[key].object
}
