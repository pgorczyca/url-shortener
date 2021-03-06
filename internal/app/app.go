package app

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/pgorczyca/url-shortener/internal/app/handler"
	"github.com/pgorczyca/url-shortener/internal/app/repository"
	"github.com/pgorczyca/url-shortener/internal/app/shortener"
	"github.com/pgorczyca/url-shortener/internal/app/utils"
	etcd "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type App struct {
	redisClient    *redis.Client
	mongoClient    *mongo.Client
	urlRepository  repository.UrlRepository
	etcdClient     *etcd.Client
	rangeProvider  shortener.RangeProvider
	shortGenerator *shortener.ShortGenerator
}

var config = utils.GetConfig()

func NewApp() (*App, error) {
	utils.InitializeLogger()
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		utils.Logger.Error("Not able to connect to mongo.", zap.Error(err))
		return nil, err
	}
	opt, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		utils.Logger.Error("Not able to connect to redis.", zap.Error(err))
		return nil, err
	}
	redisClient := redis.NewClient(opt)

	mongoRepository := repository.NewMongo(mongoClient)
	redisRepository := repository.NewRedis(redisClient, mongoRepository)

	etcdClient, err := etcd.New(etcd.Config{Endpoints: config.EtcdEndpoints, DialTimeout: 5 * time.Second})
	if err != nil {
		utils.Logger.Error("Not able to connect to etcd.", zap.Error(err))
		return nil, err
	}
	rangeProvider := shortener.NewEtcdProvider(etcdClient)
	rangeProvider.Initialize()

	shortGenerator, _ := shortener.NewShortGenerator(rangeProvider)

	return &App{
		redisClient:    redisClient,
		mongoClient:    mongoClient,
		urlRepository:  redisRepository,
		etcdClient:     etcdClient,
		rangeProvider:  rangeProvider,
		shortGenerator: shortGenerator,
	}, nil
}

func (a *App) ServeHTTP() {
	utils.Logger.Info("Running application")
	defer utils.Logger.Info("Stopping application")

	defer a.mongoClient.Disconnect(context.TODO())
	defer a.redisClient.Close()
	defer a.etcdClient.Close()
	defer utils.Logger.Sync()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/healthz", handler.Healthz)
	router.POST("/url", a.handleCreateUrlRequest(handler.CreateUrl))
	router.GET("/:short", a.handleGetUrlRequest(handler.GetUrl))
	router.Run()

}

func (a *App) handleCreateUrlRequest(handler func(c *gin.Context, repo repository.UrlRepository, sg *shortener.ShortGenerator)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c, a.urlRepository, a.shortGenerator)
	}
}

func (a *App) handleGetUrlRequest(handler func(c *gin.Context, repo repository.UrlRepository)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c, a.urlRepository)
	}
}
