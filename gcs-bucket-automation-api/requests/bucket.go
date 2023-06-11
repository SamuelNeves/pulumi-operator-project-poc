package requests

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	BucketRequest struct {
		Name                     string `json:"name" validate:"required"`
		Location                 string `json:"location" validate:"required"`
		Project                  string `json:"project" validate:"required"`
		UniformBucketLevelAccess bool   `json:"uniformBucketLevelAccess"`
		PublicAccessPrevention   string `json:"publicAccessPrevention"`
	}
	CustomValidator struct {
		Validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
