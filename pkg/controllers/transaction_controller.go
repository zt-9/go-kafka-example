package controllers

import (
	"encoding/json"
	"fmt"
	"go-kafka-example/pkg/models"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// todo

// GET `transactiones/{hash}`
// get a transaction
func GetTransaction(db *models.MetadataDB[models.Transaction]) echo.HandlerFunc {
	return func(c echo.Context) error {
		hash := c.Param("hash")
		data, err := db.Get(hash)
		if err != nil {
			if err == redis.Nil {
				return c.JSON(http.StatusNotFound, "The requested resource does not exist.")
			} else {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}

		}
		return c.JSON(http.StatusOK, data)
	}
}

// // GET `/transactiones`
// // get all transactions
func GetTransactions(db *models.MetadataDB[models.Transaction]) echo.HandlerFunc {
	return func(c echo.Context) error {
		values, err := db.GetAll(client)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, values)
	}
}

// POST `/transactiones/{hash}`
// create a transaction
func CreateTransaction(producer sarama.SyncProducer, db *models.MetadataDB[models.Transaction]) echo.HandlerFunc {
	return func(c echo.Context) error {
		hash := c.Param("hash")
		exists, err := db.KeyExists(hash)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if exists {
			return c.JSON(http.StatusBadRequest, "The resource to create already exist.")
		}

		var data models.Transaction

		// validate request body
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// use the validator to validate required field
		if validationErr := validate.Struct(&data); validationErr != nil {
			return c.JSON(http.StatusBadRequest, validationErr.Error())
		}

		// check if path param is the same as the hash field
		if hash != data.Hash {
			return c.JSON(http.StatusBadRequest, "transaction hash not match")
		}

		// create transaction
		newTx := models.Transaction{
			Hash:    data.Hash,
			Chainid: data.Chainid,
			From:    data.From,
			To:      data.To,
			Status:  data.Status,
		}

		// send kafka message
		newTxJSON, err := json.Marshal(newTx)
		if err != nil {
			return fmt.Errorf("failed to marshal transaction: %w", err)
		}

		msg := &sarama.ProducerMessage{
			Topic: "transaction",
			Key:   sarama.StringEncoder(hash),
			Value: sarama.StringEncoder(newTxJSON),
		}

		if _, _, err = producer.SendMessage(msg); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, "transaction created")

	}
}

// PUT `/transactiones/{hash}`
// update a transaction
func UpdateTransaction(producer sarama.SyncProducer, db *models.MetadataDB[models.Transaction]) echo.HandlerFunc {
	return func(c echo.Context) error {
		hash := c.Param("hash")
		exists, err := db.KeyExists(hash)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if !exists {
			return c.JSON(http.StatusBadRequest, "The resource to updare does not exist.")
		}

		var data models.Transaction
		// validate request body
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// use the validator to validate required field
		if validationErr := validate.Struct(&data); validationErr != nil {
			return c.JSON(http.StatusBadRequest, validationErr.Error())
		}

		// check if path param is the same as the hash field
		if hash != data.Hash {
			return c.JSON(http.StatusBadRequest, "transaction hash not match")
		}

		// update transaction
		newTx := models.Transaction{
			Hash:    data.Hash,
			Chainid: data.Chainid,
			From:    data.From,
			To:      data.To,
			Status:  data.Status,
		}

		// send kafka message
		newTxJSON, err := json.Marshal(newTx)
		if err != nil {
			return fmt.Errorf("failed to marshal transaction: %w", err)
		}

		msg := &sarama.ProducerMessage{
			Topic: "transaction",
			Key:   sarama.StringEncoder(hash),
			Value: sarama.StringEncoder(newTxJSON),
		}

		if _, _, err = producer.SendMessage(msg); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, "transaction updated")
	}
}

// DELETE `/transactiones/{hash}`
// delete a transaction
func DeleteTransaction(producer sarama.SyncProducer, db *models.MetadataDB[models.Transaction]) echo.HandlerFunc {
	return func(c echo.Context) error {
		hash := c.Param("hash")
		keysDeleted, err := db.Delete(hash)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if keysDeleted == 0 {
			return c.JSON(http.StatusBadRequest, "The resouce to delete does not exist.")
		}

		return c.JSON(http.StatusOK, "transaction deleted")

	}
}
