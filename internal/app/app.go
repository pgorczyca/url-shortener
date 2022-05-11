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

	router := gin.Default()
	router.GET("/healthz", handler.Healthz)
	router.POST("/url", handler.CreateUrl(a.urlRepository))
	router.Run()
	// url, esrr := redisRepository.GetByShort(context.TODO(), "ASt")
	// if esrr != nil {
	// 	fmt.Println(esrr)
	// } else {
	// 	fmt.Println(url)
	// }

	defer a.mongoClient.Disconnect(context.TODO())
	defer a.redisClient.Close()
}
