package contentdb

import (
	"github.com/tekenradar/content-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (dbService *ContentDBService) CreateIndexTickBiteMapInfos(instanceID string) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	_, err := dbService.collectionRefTickBiteMapInfos(instanceID).Indexes().CreateOne(
		ctx, mongo.IndexModel{
			Keys: bson.M{
				"time": -1,
			},
		},
	)
	return err
}

func (dbService *ContentDBService) AddTickBiteMapData(instanceID string, tickBiteMapData types.TickBiteMapData) (string, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	res, err := dbService.collectionRefTickBiteMapInfos(instanceID).InsertOne(ctx, tickBiteMapData)
	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), err
}

func (dbService *ContentDBService) FindTickBiteMapDataNewerThan(instanceID string, time int64) (tickBiteMapData []types.TickBiteMapData, err error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{}
	if time > 0 {
		filter["time"] = bson.M{"$gt": time}
	}

	batchSize := int32(32)
	opts := options.FindOptions{
		BatchSize: &batchSize,
	}
	cur, err := dbService.collectionRefTickBiteMapInfos(instanceID).Find(ctx, filter, &opts)

	if err != nil {
		return tickBiteMapData, err
	}
	defer cur.Close(ctx)

	tickBiteMapData = []types.TickBiteMapData{}
	for cur.Next(ctx) {
		var result types.TickBiteMapData
		err := cur.Decode(&result)

		if err != nil {
			return tickBiteMapData, err
		}

		tickBiteMapData = append(tickBiteMapData, result)
	}
	if err := cur.Err(); err != nil {
		return tickBiteMapData, err
	}

	return tickBiteMapData, nil
}
