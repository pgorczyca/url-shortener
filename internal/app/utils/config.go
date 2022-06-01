package utils

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AppConfig struct {
	MongoURI         string   `mapstructure:"MONGO_URI"`
	MongoDB          string   `mapstructure:"MONGO_DB"`
	RedisURL         string   `mapstructure:"REDIS_URL"`
	EtcdEndpoints    []string `mapstructure:"ETCD_ENDPOINTS"`
	CounterIncrement uint64   `mapstructure:"COUNTER_INCREMENT"`
	CounterTreshold  float64  `mapstructure:"COUNTER_TRESHOLD"`
}

var config AppConfig

func GetConfig() *AppConfig {
	return &config
}

func init() {
	vp := viper.New()
	vp.AutomaticEnv()
	vp.SetDefault("MONGO_URI", "mongodb://localhost:27017")
	vp.SetDefault("MONGO_DB", "mongoDB")
	vp.SetDefault("REDIS_URL", "redis://localhost:6379")
	vp.SetDefault("ETCD_ENDPOINTS", "localhost:2379")
	vp.SetDefault("COUNTER_INCREMENT", "100")
	vp.SetDefault("COUNTER_TRESHOLD", "0.9")

	if err := vp.Unmarshal(&config); err != nil {
		Logger.Info("Not able to unmarshall config.", zap.Error(err))

	}
}
