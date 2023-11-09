package test

import (
	"context"
	"go-kafka-example/controllers"
	"go-kafka-example/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// test for address controllers

// test for GET /addresses
func TestGetAllAaddresses(t *testing.T) {
	ctx := context.Background()
	// setup docker
	redisC, err := setUpRedis(ctx)
	if err != nil {
		t.Fatalf("failed to setup redis container: %v", err)
	}

	kafkaC, err := setUpKafka(ctx)
	if err != nil {
		t.Fatalf("failed to setup kafka container: %v", err)
	}

	// Clean up the containers after the test if complete
	t.Cleanup(func() {
		if err := redisC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate redis container: %s", err)
		}
		if err := kafkaC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate kafka container: %s", err)
		}
	})

	// set up echo
	e := echo.New()
	req := httptest.NewRequest(
		http.MethodGet, "/addresses", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db := models.NewMetadataDB[models.Address](0)

	// assertions
	if assert.NoError(t, controllers.GetAddresses(db)(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

}
