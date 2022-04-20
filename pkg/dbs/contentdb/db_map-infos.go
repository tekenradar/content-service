package contentdb

import (
	"github.com/tekenradar/content-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (dbService *ContentDBService) AddMapData(instanceID string, mapData types.MapData) (string, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	res, err := dbService.collectionRefStudyInfos(instanceID).InsertOne(ctx, mapData)
	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), err
}


func (dbService *ContentDBService) FindMapDataByTime(instanceID string, time int64) (mapData []types.MapData, err error){
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
	cur, err := dbService.collectionRefStudyInfos(instanceID).Find(ctx,filter, &opts)

	if err != nil {
		return mapData, err
	}
	defer cur.Close(ctx)

	mapData = []types.MapData{}
	for cur.Next(ctx) {
		var result types.MapData
		err := cur.Decode(&result)

		if err != nil {
			return mapData, err
		}

		mapData = append(mapData, result)
	}
	if err := cur.Err(); err != nil {
		return mapData, err
	}

	return mapData, nil
}