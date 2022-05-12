package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/pgorczyca/url-shortener/internal/app/handler"
	"github.com/pgorczyca/url-shortener/internal/app/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	redisClient   *redis.Client
	mongoClient   *mongo.Client
	urlRepository repository.UrlRepository
}

func NewApp() (*App, error) {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	opt, err := redis.ParseURL("redis://localhost:6379")
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(opt)

	mongoRepository := repository.NewMongo(mongoClient)
	redisRepository := repository.NewRedis(redisClient, mongoRepository)
	return &App{
		redisClient:   redisClient,
		mongoClient:   mongoClient,
		urlRepository: redisRepository,
	}, nil
}

func (a *App) Run() {
	defer a.mongoClient.Disconnect(context.TODO())
	defer a.redisClient.Close()

	router := gin.Default()
	router.GET("/healthz", handler.Healthz)
	router.POST("/url", a.handleGinRequest(handler.CreateUrl))
	router.GET("/url/:short", a.handleGinRequest(handler.GetUrl))
	router.Run()
}

type requestHandlerFunc func(c *gin.Context, repo repository.UrlRepository)

func (a *App) handleGinRequest(handler requestHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c, a.urlRepository)
	}
}
