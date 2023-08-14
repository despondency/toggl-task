package persister

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ResultModel struct {
	json string
}

type ResultPersister interface {
	Persist(ctx context.Context, model *ResultModel) (string, error)
}

type MongoPersister struct {
	client *mongo.Client
}

func NewMongoPersister(client *mongo.Client) ResultPersister {
	return &MongoPersister{
		client: client,
	}
}

func (mp *MongoPersister) Persist(ctx context.Context, model *ResultModel) (string, error) {
	col := mp.client.Database("toggl_database").Collection("receipt_results")
	res, err := col.InsertOne(ctx, model)
	if err != nil {
		return "", err
	}
	castedId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("could not cast insertedID to objectID")
	}
	return castedId.String(), nil
}
