package kafka

type Event struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       string
	Data      []byte
}
