package repository

import (
	"data_collector/pkg/storage"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CollectorSNG struct {
	collector chan *event
	dataMx    sync.Mutex
	data      map[string]int
	BaseCollector
}

func NewCollectorSNG(db *storage.DBConnector) *CollectorSNG {
	cl := &CollectorSNG{
		collector: make(chan *event, 1000),
		data:      make(map[string]int),
	}
	cl.db = db

	go cl.collectEvents()
	go cl.saveEvents()
	return cl
}

func (cl *CollectorSNG) HandleEvent(id int32, label string) {
	cl.collector <- &event{id: id, label: label}
}

func (cl *CollectorSNG) collectEvents() {
	for i := range cl.collector {
		lbl := strconv.Itoa(int(i.id)) + "_" + i.label
		cl.dataMx.Lock()
		if _, ok := cl.data[lbl]; !ok {
			cl.data[lbl] = 0
		}
		cl.data[lbl]++
		cl.dataMx.Unlock()
	}
}

func (cl *CollectorSNG) saveEvents() {
	for range time.Tick(5 * time.Second) {
		cl.dataMx.Lock()
		counterData := cl.data
		cl.data = make(map[string]int)
		cl.dataMx.Unlock()
		values := make([]interface{}, 0, 30)
		placeHolders := make([]string, 0, 10)
		for i, v := range counterData {
			placeHolders = append(placeHolders, "(?,?,?)")
			data := strings.Split(i, "_")
			values = append(values, data[0], data[1], v)
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

func (cl *CollectorSNG) insertData(placeHolders []string, values []interface{}) {
	placeStr := strings.Join(placeHolders, ",")
	sqlI := "INSERT INTO counters (event_id,event_label,counter) VALUES " + placeStr + " AS new(a,b,c) ON DUPLICATE KEY UPDATE counter = counter+c;"
	cl.db.Client.Exec(sqlI, values...)
}
