package persister

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	databaseName          = "toggl_db"
	receiptCollectionName = "receipt_results"
)

type ResultModel struct {
	Payload string   `json:"payload" bson:"payload"`
	Tags    []string `json:"tags" bson:"tags"`
}

type ResultPersister interface {
	Persist(ctx context.Context, model *ResultModel) (string, error)
	Get(ctx context.Context, id string) (*ResultModel, error)
	GetByTags(ctx context.Context, tags []string) ([]*ResultModel, error)
}

type MongoPersister struct {
	client *mongo.Client
}

func NewMongoPersister(client *mongo.Client) ResultPersister {
	col := client.Database(databaseName).Collection(receiptCollectionName)
	nameOfIdx, err := col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"tags": 1,
			},
		},
	)
	log.Infof("index: %s, created successfully", nameOfIdx)
	if err != nil {
		panic(err)
	}

	return &MongoPersister{
		client: client,
	}
}

func (mp *MongoPersister) Get(ctx context.Context, id string) (*ResultModel, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	col := mp.client.Database(databaseName).Collection(receiptCollectionName)
	res := &ResultModel{}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", objID}}
	err = col.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("record not found for id %s", id)
		}
	}
	return res, nil
}

func (mp *MongoPersister) GetByTags(ctx context.Context, tags []string) ([]*ResultModel, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	cursor, err := mp.client.Database("toggl_db").
		Collection("receipt_results").
		Find(ctx,
			bson.D{
				{
					"tags",
					bson.D{{"$all", tags}},
				},
			})
	if err != nil {
		return nil, err
	}
	res := make([]*ResultModel, 0)
	err = cursor.All(ctx, &res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("records with tags %s not found", tags)
		}
		return nil, err
	}
	return res, nil
}

func (mp *MongoPersister) Persist(ctx context.Context, model *ResultModel) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	col := mp.client.Database(databaseName).Collection(receiptCollectionName)
	res, err := col.InsertOne(ctx, model)
	if err != nil {
		return "", err
	}
	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("cant cast inserted id to object ID")
	}
	return objID.Hex(), nil
}
