package controllers

import (
	"go-kafka-example/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

//todo

// GET `labeles/{label}`
// get a label with all the related addresses
func GetLabel(db *models.MetadataDB[models.Label]) echo.HandlerFunc {
	return func(c echo.Context) error {
		label := c.Param("label")
		data, err := db.Get(label)
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

// GET `/labeles`
// get all labels
func GetLabels(db *models.MetadataDB[models.Label]) echo.HandlerFunc {
	return func(c echo.Context) error {
		values, err := db.GetAll(client)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, values)
	}
}
