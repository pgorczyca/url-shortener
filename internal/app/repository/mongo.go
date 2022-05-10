package repository

import (
	"context"
	"errors"

	"github.com/pgorczyca/url-shortener/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	mongoDB        = "mongoDB"
	mongoColletion = "urls"
)

type MongoUrlRepository struct {
	collection *mongo.Collection
}

func NewMongo(client *mongo.Client) *MongoUrlRepository {
	return &MongoUrlRepository{collection: client.Database(mongoDB).Collection(mongoColletion)}
}
func (r *MongoUrlRepository) Add(ctx context.Context, u model.Url) error {

	doc := bson.M{"long": u.Long, "short": u.Short, "expired_at": u.ExpiredAt, "created_at": u.CreatedAt}
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}
func (r *MongoUrlRepository) GetByShort(ctx context.Context, short string) (model.Url, error) {
	var u bson.M
	if err := r.collection.FindOne(ctx, bson.M{"short": short}).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Url{}, errors.New("no results")
		}
		return model.Url{}, err
	}
	return model.Url{
		Long:      u["long"].(string),
		Short:     u["short"].(string),
		ExpiredAt: u["expired_at"].(primitive.DateTime).Time(),
		CreatedAt: u["created_at"].(primitive.DateTime).Time(),
	}, nil
}
