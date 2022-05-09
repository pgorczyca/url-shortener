package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://localhost:27017"

type url struct {
	Long      string    `json:"long"`
	Short     string    `json:"short"`
	ExpiredAt time.Time `json:"exipred_at"`
	CreatedAt time.Time `json:"created_at"`
}
type urlRepository interface {
	add(ctx context.Context, url url) error
	getByShort(ctx context.Context, short string) (url, error)
}

type mongoUrlRepository struct {
	collection *mongo.Collection
}

func (r *mongoUrlRepository) add(ctx context.Context, u url) error {

	doc := bson.M{"long": u.Long, "short": u.Short, "expired_at": u.ExpiredAt, "created_at": u.CreatedAt}
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}
func (r *mongoUrlRepository) getByShort(ctx context.Context, short string) (url, error) {
	var u bson.M
	if err := r.collection.FindOne(ctx, bson.M{"short": short}).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return url{}, errors.New("no results")
		}
		return url{}, err
	}
	return url{
		Long:      u["long"].(string),
		Short:     u["short"].(string),
		ExpiredAt: u["expired_at"].(primitive.DateTime).Time(),
		CreatedAt: u["created_at"].(primitive.DateTime).Time(),
	}, nil
}

type redisUrlRepository struct {
	client  *redis.Client
	urlRepo urlRepository
}

func (r *redisUrlRepository) add(ctx context.Context, u url) error {
	r.urlRepo.add(ctx, u)
	jsonUrl, err := json.Marshal(u)
	if err != nil {
		return err
	}
	r.client.Set(u.Short, jsonUrl, 0)
	return nil
}

func (r *redisUrlRepository) getByShort(ctx context.Context, short string) (url, error) {
	jsonUrl := r.client.Get(short)
	_, err := jsonUrl.Result()
	if err == redis.Nil {
		return r.urlRepo.getByShort(ctx, short)
	}
	var u url
	bytes, err := jsonUrl.Bytes()
	if err != nil {
		return url{}, err
	}
	json.Unmarshal(bytes, &u)
	return u, nil
}
func main() {
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("myDB").Collection("urls")
	mongoRepository := mongoUrlRepository{collection: coll}

	opt, err := redis.ParseURL("redis://localhost:6379")
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(opt)
	repo := redisUrlRepository{client: redisClient, urlRepo: &mongoRepository}
	// url1 := url{
	// 	Long:      "https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/insert/",
	// 	Short:     "ASt",
	// 	ExpiredAt: time.Now().Add(time.Hour * 6),
	// 	CreatedAt: time.Now(),
	// }
	// repo.add(context.TODO(), url1)
	url, esrr := repo.getByShort(context.TODO(), "ASt")
	if esrr != nil {
		fmt.Println(esrr)
	} else {
		fmt.Println(url)
	}
}
