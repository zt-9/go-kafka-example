package controllers

import (
	"encoding/json"
	"fmt"
	"go-kafka-example/config"
	"go-kafka-example/models"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

var client = redis.NewClient(&redis.Options{
	Addr:     config.RedisServerAddr,
	Password: "", // no password set
	DB:       0,  // use default DB
})

// GET `addresses/{address}`
// get an address
func GetAddress(db *models.MetadataDB[models.Address]) echo.HandlerFunc {
	return func(c echo.Context) error {
		address := c.Param("address")
		data, err := db.Get(address)
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

// GET `/addresses`
// get all addresses
func GetAddresses(db *models.MetadataDB[models.Address]) echo.HandlerFunc {
	return func(c echo.Context) error {
		values, err := db.GetAll(client)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, values)
	}
}

var validate = validator.New()

// POST `/addresses/{address}`
// create an address
func CreateAddress(producer sarama.SyncProducer, db *models.MetadataDB[models.Address]) echo.HandlerFunc {
	return func(c echo.Context) error {
		address := c.Param("address")
		exists, err := db.KeyExists(address)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if exists {
			return c.JSON(http.StatusBadRequest, "The resource to create already exist.")
		}

		var data models.Address

		// validate the request body
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// use the validator to validate required field
		if validationErr := validate.Struct(&data); validationErr != nil {
			return c.JSON(http.StatusBadRequest, validationErr.Error())
		}

		// check if path param is the same as the address field
		if address != data.Address {
			return c.JSON(http.StatusBadRequest, "address not match")
		}

		// create address
		newAddr := models.Address{
			Address: data.Address,
			Labels:  data.Labels,
		}

		// send kafka message
		newAddrJSON, err := json.Marshal(newAddr)
		if err != nil {
			return fmt.Errorf("failed to marshal Address: %w", err)
		}

		msg := &sarama.ProducerMessage{
			Topic: "address",
			Key:   sarama.StringEncoder(address),
			Value: sarama.StringEncoder(newAddrJSON),
		}

		if _, _, err = producer.SendMessage(msg); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, "address created")

	}
}

// PUT `/addresses/{address}`
// update an address
func UpdateAddress(producer sarama.SyncProducer, db *models.MetadataDB[models.Address]) echo.HandlerFunc {
	return func(c echo.Context) error {
		address := c.Param("address")
		exists, err := db.KeyExists(address)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if !exists {
			return c.JSON(http.StatusBadRequest, "The resource to updare does not exist.")
		}

		var data models.Address

		// validate the request body
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// use the validator to validate required field
		if validationErr := validate.Struct(&data); validationErr != nil {
			return c.JSON(http.StatusBadRequest, validationErr.Error())
		}

		// check if path param is the same as the address field
		if address != data.Address {
			return c.JSON(http.StatusBadRequest, "address not match")
		}

		// create address
		newAddr := models.Address{
			Address: data.Address,
			Labels:  data.Labels,
		}

		// send kafka message
		newAddrJSON, err := json.Marshal(newAddr)
		if err != nil {
			return fmt.Errorf("failed to marshal Address: %w", err)
		}

		msg := &sarama.ProducerMessage{
			Topic: "address",
			Key:   sarama.StringEncoder(address),
			Value: sarama.StringEncoder(newAddrJSON),
		}

		if _, _, err = producer.SendMessage(msg); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, "address updated")

	}
}

// DELETE `/addresses/{address}`
// delete an address
func DeleteAddress(producer sarama.SyncProducer, db *models.MetadataDB[models.Address]) echo.HandlerFunc {
	return func(c echo.Context) error {
		address := c.Param("address")
		keysDeleted, err := db.Delete(address)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if keysDeleted == 0 {
			return c.JSON(http.StatusBadRequest, "The resouce to delete does not exist.")
		}

		return c.JSON(http.StatusOK, "address deleted")

	}
}
