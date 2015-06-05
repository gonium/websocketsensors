package websocketsensors

const (
	Celsius = iota
	Watt    = iota
)

const (
	STATUS_OK              = iota
	STATUS_REGISTER_FAILED = iota
)

type RegisterMessage struct {
	Unit int
}

type UpdateMessage struct {
	Timestamp int
	Value     float32
}

type UpdateResponse struct {
	Status int
}
