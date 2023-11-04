package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-kafka-example/config"
	"go-kafka-example/pkg/controllers"
	"go-kafka-example/pkg/models"
	"go-kafka-example/utils"
	"log"

	"github.com/IBM/sarama"
	"github.com/labstack/echo/v4"
)

// Implements sarama.ConsumerGroupHandler.
type ConsumerGroupHandler struct {
	addressDB     *models.MetadataDB[models.Address]
	labelDB       *models.MetadataDB[models.Label]
	transactionDB *models.MetadataDB[models.Transaction]
}

func (*ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil } // required to satify the sarama.ConsumerGroupHandler interface
func (*ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (handler *ConsumerGroupHandler) ConsumeClaim(
	sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		topic := msg.Topic
		key := string(msg.Key)

		log.Printf("Consuming topic %s", topic)

		switch topic {
		case "address":
			var data models.Address
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Printf("failed to unmarshal address: %v", err)
			}

			if err := handler.addressDB.Update(key, data); err != nil {
				log.Printf("failed to update db: %v", err)
			}
			// add to new labels to label store
			if err := utils.StoreAddressLabels(key, data.Labels, handler.labelDB); err != nil {
				log.Printf("failed to  store labels: %v", err)
			}
			sess.MarkMessage(msg, "") // mark the msg as processed
		case "label":
			var data models.Label
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Printf("failed to unmarshal address: %v", err)
			}

			if err := handler.labelDB.Update(key, data); err != nil {
				log.Printf("failed to update db: %v", err)
			}
			sess.MarkMessage(msg, "") // mark the msg as processed
		case "transaction":
			var data models.Transaction
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Printf("failed to unmarshal address: %v", err)
			}

			if err := handler.transactionDB.Update(key, data); err != nil {
				log.Printf("failed to update db: %v", err)
			}
			sess.MarkMessage(msg, "") // mark the msg as processed
		}
	}
	return nil
}

func initConsumerGroup() (sarama.ConsumerGroup, error) {
	saramaConfig := sarama.NewConfig()
	consumerGroup, err := sarama.NewConsumerGroup([]string{config.KafkaServerAddr}, config.ConsumerGroup, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize consumer group: %w", err)
	}
	return consumerGroup, nil
}

func setUpConsumerGroup(ctx context.Context, addressDB *models.MetadataDB[models.Address], labelDB *models.MetadataDB[models.Label], transactionDB *models.MetadataDB[models.Transaction]) {
	consumerGroup, err := initConsumerGroup()
	if err != nil {
		log.Printf("consumer initialization error: %v", err)
	}

	defer consumerGroup.Close()

	consumer := &ConsumerGroupHandler{
		addressDB:     addressDB,
		labelDB:       labelDB,
		transactionDB: transactionDB,
	}

	for { // run for loop indefinitely. consume messages from the topic and processing any errors that arise
		if err = consumerGroup.Consume(ctx, config.Topics[:], consumer); err != nil {
			log.Printf("error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			return
		}

	}
}

func main() {
	// setup a cancellable events that can be used to stop the consumer group
	ctx, cancel := context.WithCancel(context.Background())

	// start consumer group in separate goroutine
	go setUpConsumerGroup(ctx, models.AddressDB, models.LabelDB, models.TransactionDB)
	defer cancel()
	fmt.Printf("Kafka CONSUMER (Group: %s) 👥📥 "+
		"started at http://localhost%s\n", config.ConsumerGroup, config.ConsumerPort)

	// routes
	e := echo.New()
	e.GET("/addresses/:address", controllers.GetAddress(models.AddressDB))
	e.GET("/addresses", controllers.GetAddresses(models.AddressDB))

	// LabelRoutes(e)
	e.GET("/labels/:label", controllers.GetLabel(models.LabelDB))
	e.GET("/labels", controllers.GetLabels(models.LabelDB))
	// TransactionRoutes(e)
	e.GET("/transactions/:hash", controllers.GetTransaction(models.TransactionDB))
	e.GET("/transactions", controllers.GetTransactions(models.TransactionDB))

	e.Logger.Fatal(e.Start(config.ConsumerPort))

}
