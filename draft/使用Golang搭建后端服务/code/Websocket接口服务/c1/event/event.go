package event

type Event struct {
	EventType string `json:"name"`
	Message   string `json:"message"`
}
