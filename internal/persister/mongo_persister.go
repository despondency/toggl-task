package persister

import (
	"context"
	"errors"
	"fmt"
	"github.com/despondency/toggl-task/internal/model"
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

type ResultPersister interface {
	Persist(ctx context.Context, model *model.Receipt) (*model.Receipt, error)
	Get(ctx context.Context, id string) (*model.Receipt, error)
	GetByTags(ctx context.Context, tags []string) ([]*model.Receipt, error)
}

type MongoPersister struct {
	client *mongo.Client
}

func NewMongoPersister(client *mongo.Client) ResultPersister {
	col := client.Database(databaseName).Collection(receiptCollectionName)
	// index creation is idempotent in mongo, so we can do it on startup always
	// for older collections it might take a while so this need to be taken in consideration!
	nameOfIdx, err := col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{
					Key:   "tags",
					Value: 1,
				},
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

func (mp *MongoPersister) Get(ctx context.Context, id string) (*model.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	col := mp.client.Database(databaseName).Collection(receiptCollectionName)
	res := &model.Receipt{}
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

func (mp *MongoPersister) GetByTags(ctx context.Context, tags []string) ([]*model.Receipt, error) {
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
	res := make([]*model.Receipt, 0)
	err = cursor.All(ctx, &res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("records with tags %s not found", tags)
		}
		return nil, err
	}
	return res, nil
}

func (mp *MongoPersister) Persist(ctx context.Context, model *model.Receipt) (*model.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	col := mp.client.Database(databaseName).Collection(receiptCollectionName)
	model.Id = primitive.NewObjectID()
	res, err := col.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}
	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("cant cast inserted id to object ID")
	}
	return model, nil
}
