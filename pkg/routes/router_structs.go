package routes

type PayloadMessage struct {
	ID    int32  `json:"id"`
	Label string `json:"label"`
}

type ICollector interface {
	HandleEvent(id int32, label string) bool
}
