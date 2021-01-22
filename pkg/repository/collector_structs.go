package repository

type event struct {
	id    int32
	label string
}

type EventStat struct {
	EventLabel string `json:"event_label" db:"event_label"`
	EventID    int32  `json:"event_id" db:"event_id"`
	Counter    int32  `json:"counter" db:"counter"`
}
