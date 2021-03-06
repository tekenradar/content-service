package contentdb

import (
	"github.com/tekenradar/content-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (dbService *ContentDBService) CreateIndexNewsItemInfos(instanceID string) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	_, err := dbService.collectionRefNewsItems(instanceID).Indexes().CreateOne(
		ctx, mongo.IndexModel{
			Keys: bson.M{
				"time": -1,
			},
		},
	)
	return err
}

func (dbService *ContentDBService) AddNewsItem(instanceID string, newsItem types.NewsItem) (string, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	res, err := dbService.collectionRefNewsItems(instanceID).InsertOne(ctx, newsItem)
	if err != nil {
		return "", err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), err
}

func (dbService *ContentDBService) GetNewsItemsList(instanceID string) (newsItemList []types.NewsItem, err error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{}
	batchSize := int32(32)
	opts := options.FindOptions{
		BatchSize: &batchSize,
	}
	cur, err := dbService.collectionRefNewsItems(instanceID).Find(ctx, filter, &opts)

	if err != nil {
		return newsItemList, err
	}
	defer cur.Close(ctx)

	newsItemList = []types.NewsItem{}
	for cur.Next(ctx) {
		var result types.NewsItem

		err := cur.Decode(&result)
		if err != nil {
			return newsItemList, err
		}

		newsItemList = append(newsItemList, result)
	}
	if err := cur.Err(); err != nil {
		return newsItemList, err
	}

	return newsItemList, nil
}
