package repository

import (
	"context"

	"github.com/pgorczyca/url-shortener/internal/app/model"
)

type UrlRepository interface {
	Add(ctx context.Context, url model.Url) error
	GetByShort(ctx context.Context, short string) (model.Url, error)
}
