package main

import (
	"automation-api/docs"
	_ "automation-api/docs"
	"automation-api/handler"
	"automation-api/requests"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"
)

var project = "bucket-automation-api"

// @title BucketRequest Factory Automation Api
// @version 1.0
// @description This api is responsible to generate buckets in a rest way

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 0.0.0.0:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /
// @schemes http,https
func main() {
	// Echo instance
	e := echo.New()
	e.Validator = &requests.CustomValidator{Validator: validator.New()}

	docs.SwaggerInfo.Host = os.Getenv("BASE_HOST")
	e.GET("/", HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	bucketHandler := handler.NewBucketHandler(project)
	e.GET("/bucket", bucketHandler.List)
	e.GET("/bucket/:name", bucketHandler.Get)
	e.POST("/bucket", bucketHandler.Post)
	e.DELETE("/bucket/:name", bucketHandler.Delete)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
