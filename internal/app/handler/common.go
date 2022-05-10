package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	apphttp "github.com/pgorczyca/url-shortener/internal/app/http"
)

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

var valid = validator.New()

func validate(r apphttp.CreateUrlRequest) ([]*ValidationError, error) {
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
