package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/pgorczyca/url-shortener/internal/app/handler"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	redisClient *redis.Client
	mongoClient *mongo.Client
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
	return &App{redisClient: redisClient, mongoClient: mongoClient}, nil

}

func (a *App) Run() {
	// mongoRepository := repository.NewMongo(a.mongoClient)
	// redisRepository := repository.NewRedis(a.redisClient, mongoRepository)

	router := gin.Default()
	router.GET("/healthz", handler.Healthz)
	router.POST("/url", handler.CreateUrl)
	router.Run()
	// url, esrr := redisRepository.GetByShort(context.TODO(), "ASt")
	// if esrr != nil {
	// 	fmt.Println(esrr)
	// } else {
	// 	fmt.Println(url)
	// }

	// url1 := model.Url{
	// 	Long:      "https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/insert/",
	// 	Short:     "ASt",
	// 	ExpiredAt: time.Now().Add(time.Hour * 6),
	// 	CreatedAt: time.Now(),
	// }
	// redisRepository.Add(context.TODO(), url1)
	defer a.mongoClient.Disconnect(context.TODO())
	defer a.redisClient.Close()
}
