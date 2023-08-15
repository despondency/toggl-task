package persister

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ResultModel struct {
	UUID    string `json:"uuid" bson:"uuid"`
	Payload string `json:"payload" bson:"payload"`
}

type ResultPersister interface {
	Persist(ctx context.Context, model *ResultModel) (string, error)
	Get(ctx context.Context, id string) (*ResultModel, error)
}

type MongoPersister struct {
	client *mongo.Client
}

func NewMongoPersister(client *mongo.Client) ResultPersister {
	return &MongoPersister{
		client: client,
	}
}

func (mp *MongoPersister) Get(ctx context.Context, id string) (*ResultModel, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	col := mp.client.Database("mongo").Collection("receipt_results")
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

func (mp *MongoPersister) Persist(ctx context.Context, model *ResultModel) (string, error) {
	col := mp.client.Database("mongo").Collection("receipt_results")
	res, err := col.InsertOne(ctx, model)
	model.UUID = uuid.New().String()
	if err != nil {
		return "", err
	}
	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("cant cast inserted id to object ID")
	}
	return objID.String(), nil
}
