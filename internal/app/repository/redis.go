package repository

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/pgorczyca/url-shortener/internal/app/model"
)

type RedisUrlRepository struct {
	client *redis.Client
	repo   UrlRepository
}

func NewRedis(redisClient *redis.Client, repo UrlRepository) *RedisUrlRepository {
	return &RedisUrlRepository{client: redisClient, repo: repo}
}

func (r *RedisUrlRepository) Add(ctx context.Context, u model.Url) error {
	r.repo.Add(ctx, u)
	jsonUrl, err := json.Marshal(u)
	if err != nil {
		return err
	}
	r.client.Set(u.Short, jsonUrl, 0)
	return nil
}

func (r *RedisUrlRepository) GetByShort(ctx context.Context, short string) (model.Url, error) {
	jsonUrl := r.client.Get(short)
	_, err := jsonUrl.Result()
	if err == redis.Nil {
		return r.repo.GetByShort(ctx, short)
	}
	var u model.Url
	bytes, err := jsonUrl.Bytes()
	if err != nil {
		return model.Url{}, err
	}
	json.Unmarshal(bytes, &u)
	return u, nil
}
