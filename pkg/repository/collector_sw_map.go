package repository

import (
	"data_collector/pkg/storage"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const rangeMaxSw int32 = 4

type CollectorSW struct {
	dataMxContainer    []*sync.Mutex
	dataContainer      []sync.Map
	collectorContainer []chan *event
	reqCounter         int32
	BaseCollector
}

func NewCollectorSW(db *storage.DBConnector) *CollectorSW {
	cl := &CollectorSW{
		dataContainer:      make([]sync.Map, rangeMaxSw, rangeMaxSw),
		collectorContainer: make([]chan *event, rangeMaxSw, rangeMaxSw),
		dataMxContainer:    make([]*sync.Mutex, rangeMaxSw, rangeMaxSw),
		reqCounter:         0,
	}
	cl.db = db

	for i := 0; i < int(rangeMax); i++ {
		//cl.dataContainer[i] = make(map[string]int, 1000)
		cl.dataContainer[i] = sync.Map{}
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

func (cl *CollectorSW) HandleEvent(id int32, label string) {
	i := atomic.AddInt32(&cl.reqCounter, 1) % rangeMax
	cl.collectorContainer[i] <- &event{id: id, label: label}
}

func (cl *CollectorSW) collectEvents(i int) {
	for e := range cl.collectorContainer[i] {
		lbl := strconv.Itoa(int(e.id)) + "_" + e.label
		v, ok := cl.dataContainer[i].Load(lbl)
		val := 1
		if ok {
			val += v.(int)
		}
		cl.dataContainer[i].LoadOrStore(lbl, val)
	}
}

func (cl *CollectorSW) saveEvents(i int) {
	for range time.Tick(5 * time.Second) {
		cl.dataMxContainer[i].Lock()
		counterData := cl.dataContainer[i]
		cl.dataContainer[i] = sync.Map{} //make(map[string]int, len(counterData))
		cl.dataMxContainer[i].Unlock()
		values := make([]interface{}, 0, 30)
		placeHolders := make([]string, 0, 10)
		counterData.Range(func(j, v interface{}) bool {
			//for j, v := range counterData {
			placeHolders = append(placeHolders, "(?,?,?)")
			data := strings.Split(j.(string), "_")
			values = append(values, data[0], data[1], v)
			if len(placeHolders) >= 10 {
				cl.insertData(placeHolders, values)
				values = make([]interface{}, 0, 30)
				placeHolders = make([]string, 0, 10)
			}
			return true
		})
		if len(placeHolders) > 0 {
			cl.insertData(placeHolders, values)
		}
	}
}
