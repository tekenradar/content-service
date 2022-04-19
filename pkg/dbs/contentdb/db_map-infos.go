package contentdb

import (
	"github.com/tekenradar/content-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (dbService *ContentDBService) AddMapData(instanceID string, mapData types.MapData) (string, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	res, err := dbService.collectionRefStudyInfos(instanceID).InsertOne(ctx, mapData)
	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), err
}