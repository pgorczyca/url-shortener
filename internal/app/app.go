package app

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/pgorczyca/url-shortener/internal/app/counter"
	"github.com/pgorczyca/url-shortener/internal/app/repository"
	etcd "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	redisClient     *redis.Client
	mongoClient     *mongo.Client
	urlRepository   repository.UrlRepository
	etcdClient      *etcd.Client
	counterProvider counter.RangeProvider
	counterManager  *counter.ShortGenerator
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

	etcdClient, err := etcd.New(etcd.Config{Endpoints: []string{"localhost:2379"}, DialTimeout: 5 * time.Second})
	if err != nil {
		return nil, err
	}
	counterProvider := counter.NewEtcdProvider(etcdClient)

	counterManager, _ := counter.NewCounterManager(counterProvider)

	return &App{
		redisClient:     redisClient,
		mongoClient:     mongoClient,
		urlRepository:   redisRepository,
		etcdClient:      etcdClient,
		counterProvider: counterProvider,
		counterManager:  counterManager,
	}, nil
}

func (a *App) Run() {
	defer a.mongoClient.Disconnect(context.TODO())
	defer a.redisClient.Close()
	defer a.etcdClient.Close()
	cm := a.counterManager

	var wg sync.WaitGroup

	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cm.GetShort()
		}()
	}
	wg.Wait()
	// router := gin.Default()
	// router.GET("/healthz", handler.Healthz)
	// router.POST("/url", a.handleGinRequest(handler.CreateUrl))
	// router.GET("/url/:short", a.handleGinRequest(handler.GetUrl))
	// router.Run()

}

type requestHandlerFunc func(c *gin.Context, repo repository.UrlRepository)

func (a *App) handleGinRequest(handler requestHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c, a.urlRepository)
	}
}
