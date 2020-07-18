package mongodb

import (
	"context"

	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

const (
	dbName = "wiki"
	colName = "articles"
)

type Storage struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func NewStorage(connStr string) (*Storage, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	return &Storage{
		client: client,
		db: client.Database(dbName),
		collection: client.Database(dbName).Collection(colName),
	}, err
}

func (s *Storage) InsertArticle(article Article) error {
	article.ID = bson.NewObjectId()
	_, err := s.collection.InsertOne(context.TODO(), article)
	return err
}