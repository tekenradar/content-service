package contentdb

import (
	"github.com/tekenradar/content-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (dbService *ContentDBService) GetNewsItemsList(instanceID string) (NewsItemList []types.NewsItem, err error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{}
	batchSize := int32(32)
	opts := options.FindOptions{
		BatchSize: &batchSize,
	}
	cur, err := dbService.collectionRefNewsItems(instanceID).Find(ctx, filter, &opts)

	if err != nil {
		return NewsItemList, err
	}
	defer cur.Close(ctx)

	NewsItemList = []types.NewsItem{}
	for cur.Next(ctx) {
		var result types.NewsItem
		err := cur.Decode(&result)

		if err != nil {
			return NewsItemList, err
		}

		NewsItemList = append(NewsItemList, result)
	}
	if err := cur.Err(); err != nil {
		return NewsItemList, err
	}

	return NewsItemList, nil
}
