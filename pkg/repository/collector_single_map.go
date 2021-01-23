package repository

import (
	"sync"
	"time"
)

type CollectorSNG struct {
	collector chan *event
	dataMx    sync.Mutex
	data      map[int32]map[string]int
	BaseCollector
}

func NewCollectorSNG(db IDatabase) *CollectorSNG {
	cl := &CollectorSNG{
		collector: make(chan *event, 1000),
		data:      make(map[int32]map[string]int),
	}
	cl.db = db

	go cl.collectEventsSNG()
	go cl.saveEvents()
	return cl
}

func (cl *CollectorSNG) HandleEvent(id int32, label string) {
	cl.collector <- &event{id: id, label: label}
}

func (cl *CollectorSNG) collectEventsSNG() {
	for i := range cl.collector {
		//lbl := strconv.Itoa(int(i.id)) + "_" + i.label
		cl.dataMx.Lock()
		if _, ok := cl.data[i.id]; !ok {
			cl.data[i.id] = make(map[string]int)
		}
		if _, ok := cl.data[i.id][i.label]; !ok {
			cl.data[i.id][i.label] = 0
		}
		cl.data[i.id][i.label]++
		cl.dataMx.Unlock()
	}
}

func (cl *CollectorSNG) saveEvents() {
	for range time.Tick(5 * time.Second) {
		cl.dataMx.Lock()
		counterData := cl.data
		cl.data = make(map[int32]map[string]int, len(counterData))
		cl.dataMx.Unlock()
		values := make([]interface{}, 0, 30)
		placeHolders := make([]string, 0, 10)
		for j, v := range counterData {
			for l, c := range v {
				placeHolders = append(placeHolders, "(?,?,?)")
				values = append(values, j, l, c)
			}
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
