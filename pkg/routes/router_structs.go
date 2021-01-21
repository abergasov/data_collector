package routes

type ICollector interface {
	HandleEvent(id int32, label string)
}
