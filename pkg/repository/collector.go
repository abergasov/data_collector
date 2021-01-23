package repository

import (
	"sync"
	"sync/atomic"
	"time"
)

const rangeMax int32 = 4

type Collector struct {
	dataMxContainer    []*sync.Mutex
	dataContainer      []map[int32]map[string]int
	collectorContainer []chan *event
	reqCounter         int32
	BaseCollector
}

func NewCollector(db IDatabase) *Collector {
	cl := &Collector{
		dataContainer:      make([]map[int32]map[string]int, rangeMax, rangeMax),
		collectorContainer: make([]chan *event, rangeMax, rangeMax),
		dataMxContainer:    make([]*sync.Mutex, rangeMax, rangeMax),
		reqCounter:         0,
	}
	cl.db = db

	for i := 0; i < int(rangeMax); i++ {
		cl.dataContainer[i] = make(map[int32]map[string]int, 1000)
		cl.collectorContainer[i] = make(chan *event, 1000)
		cl.dataMxContainer[i] = &sync.Mutex{}
		go cl.collectEvents(i)
		go cl.saveEvents(i)
	}

	go func() {
		for range time.Tick(1 * time.Second) {
			atomic.StoreInt32(&cl.reqCounter, 0)
		}
	}()

	return cl
}

func (cl *Collector) HandleEvent(id int32, label string) {
	i := atomic.AddInt32(&cl.reqCounter, 1) % rangeMax
	cl.collectorContainer[i] <- &event{id: id, label: label}
}

func (cl *Collector) collectEvents(i int) {
	for e := range cl.collectorContainer[i] {
		cl.dataMxContainer[i].Lock()
		if _, ok := cl.dataContainer[i][e.id]; !ok {
			cl.dataContainer[i][e.id] = make(map[string]int)
		}
		if _, ok := cl.dataContainer[i][e.id][e.label]; !ok {
			cl.dataContainer[i][e.id][e.label] = 0
		}
		cl.dataContainer[i][e.id][e.label]++
		cl.dataMxContainer[i].Unlock()
	}
}

func (cl *Collector) saveEvents(i int) {
	for range time.Tick(5 * time.Second) {
		cl.dataMxContainer[i].Lock()
		counterData := cl.dataContainer[i]
		cl.dataContainer[i] = make(map[int32]map[string]int, len(counterData))
		cl.dataMxContainer[i].Unlock()
		values := make([]interface{}, 0, 30)
		placeHolders := make([]string, 0, 10)
		for j, v := range counterData {
			for l, c := range v {
				placeHolders = append(placeHolders, "(?,?,?)")
				values = append(values, j, l, c)
			}
			//data := strings.Split(j, "_")
			if len(placeHolders) >= 10 {
				cl.insertData(placeHolders, values)
				values = make([]interface{}, 0, 30)
				placeHolders = make([]string, 0, 10)
			}
		}
		if len(placeHolders) > 0 {
			cl.insertData(placeHolders, values)
		}
	}
}
