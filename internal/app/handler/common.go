package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	apphttp "github.com/pgorczyca/url-shortener/internal/app/http"
)

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

var valid = validator.New()

func validateCreate(r apphttp.CreateUrlRequest) ([]*ValidationError, error) {
	var errs []*ValidationError
	err := valid.Struct(r)
	if err == nil {
		return errs, nil
	}

	for _, err := range err.(validator.ValidationErrors) {
		er := ValidationError{Field: err.Field(), Tag: err.Tag(), Value: err.Error()}
		errs = append(errs, &er)
	}

	return errs, errors.New("validation error")
}

func internalServerErrorResponse(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal server error",
	})
}
