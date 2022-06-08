package contentdb

import (
	"github.com/tekenradar/content-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (dbService *ContentDBService) SaveFileInfo(instanceID string, fileInfo types.FileInfo) (types.FileInfo, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	if fileInfo.ID.IsZero() {
		fileInfo.ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": fileInfo.ID}

	upsert := true
	rd := options.After
	options := options.FindOneAndReplaceOptions{
		Upsert:         &upsert,
		ReturnDocument: &rd,
	}
	elem := types.FileInfo{}
	err := dbService.collectionRefUploadedFiles(instanceID).FindOneAndReplace(
		ctx, filter, fileInfo, &options,
	).Decode(&elem)
	return elem, err
}
