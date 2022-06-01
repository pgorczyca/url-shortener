package app

type AppConfig struct {
	Environment      string
	MongoURI         string
	MongoDB          string
	RedisURL         string
	EtcdEndpoints    []string
	CounterIncrement uint64
	CounterTreshold  float64
}

func loadConfig() *AppConfig {

	return &AppConfig{}
}
