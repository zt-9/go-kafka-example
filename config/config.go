package config

const (
	ConsumerGroup   = "metadata-group"
	ConsumerPort    = ":8081"
	ProducerPort    = ":8080"
	KafkaServerAddr = "localhost:9092"
	RedisServerAddr = "localhost:6379"
)

var Topics = [3]string{"address", "label", "transaction"}
