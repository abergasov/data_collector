package repository

type Collector struct {
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) HandleEvent(id int32, label string) {

}
