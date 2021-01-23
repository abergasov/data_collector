package repository

import (
	"data_collector/pkg/logger"
	"database/sql"
	"strings"
)

type event struct {
	id    int32
	label string
}

type EventStat struct {
	EventLabel string `json:"event_label" db:"event_label"`
	EventID    int32  `json:"event_id" db:"event_id"`
	Counter    int32  `json:"counter" db:"counter"`
}

//go:generate mockgen -source=collector_structs.go -destination=collector_structs_mock.go -package=repository
type IDatabase interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type BaseCollector struct {
	db IDatabase
}

func (cl *BaseCollector) insertData(placeHolders []string, values []interface{}) {
	placeStr := strings.Join(placeHolders, ",")
	sqlI := "INSERT INTO counters (event_id,event_label,counter) VALUES " + placeStr + " AS new(a,b,c) ON DUPLICATE KEY UPDATE counter = counter+c;"
	cl.db.Exec(sqlI, values...)
}

func (cl *BaseCollector) GetState() []EventStat {
	var p []EventStat
	err := cl.db.Select(&p, "SELECT event_id, event_label, counter FROM counters LIMIT 20")
	if err != nil {
		logger.Error("error load stat", err)
	}
	return p
}
