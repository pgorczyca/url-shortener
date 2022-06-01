package repository

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/pgorczyca/url-shortener/internal/app/model"
	"github.com/pgorczyca/url-shortener/internal/app/utils"
	"go.uber.org/zap"
)

type RedisUrlRepository struct {
	client *redis.Client
	repo   UrlRepository
}

func NewRedis(redisClient *redis.Client, repo UrlRepository) *RedisUrlRepository {
	return &RedisUrlRepository{client: redisClient, repo: repo}
}

func (r *RedisUrlRepository) Add(ctx context.Context, u model.Url) error {
	err := r.repo.Add(ctx, u)
	if err != nil {
		utils.Logger.Error("Not able to insert to repository.", zap.Error(err))
		return err
	}
	jsonUrl, err := json.Marshal(u)
	if err != nil {
		utils.Logger.Error("Not able to marshal json.", zap.Error(err))
		return err
	}

	redistStatusCmd := r.client.Set(u.Short, jsonUrl, 0)
	if redistStatusCmd.Err() != nil {
		utils.Logger.Error("Not able to save to redis.", zap.Error(err))
		return redistStatusCmd.Err()
	}
	return nil
}

func (r *RedisUrlRepository) GetByShort(ctx context.Context, short string) (model.Url, error) {
	jsonUrl := r.client.Get(short)
	_, err := jsonUrl.Result()
	if err == redis.Nil {
		utils.Logger.Info("Getting record from another repository.")
		return r.repo.GetByShort(ctx, short)
	}
	bytes, err := jsonUrl.Bytes()
	if err != nil {
		utils.Logger.Error("Not able to read bytes from redis.", zap.Error(err))
		return model.Url{}, err
	}
	var u model.Url
	json.Unmarshal(bytes, &u)
	return u, nil
}
