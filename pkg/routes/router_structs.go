package routes

import "data_collector/pkg/repository"

type PayloadMessage struct {
	ID    int32  `json:"id"`
	Label string `json:"label"`
}

type ICollector interface {
	HandleEvent(id int32, label string)
	GetState() []repository.EventStat
}
