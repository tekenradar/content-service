package contentdb

import (
	"errors"

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

func (dbService *ContentDBService) FindFileInfo(instanceID string, fileID string) (types.FileInfo, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(fileID)
	filter := bson.M{"_id": _id}

	elem := types.FileInfo{}
	err := dbService.collectionRefUploadedFiles(instanceID).FindOne(ctx, filter).Decode(&elem)
	return elem, err
}

func (dbService *ContentDBService) DeleteFileInfo(instanceID string, fileID string) (count int64, err error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	if fileID == "" {
		return 0, errors.New("file id must be defined")
	}
	_id, _ := primitive.ObjectIDFromHex(fileID)
	filter := bson.M{"_id": _id}

	res, err := dbService.collectionRefUploadedFiles(instanceID).DeleteOne(ctx, filter)
	return res.DeletedCount, err
}

func (dbService *ContentDBService) GetFileInfoList(instanceID string) (fileInfoList []types.FileInfo, err error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{}
	batchSize := int32(32)
	opts := options.FindOptions{
		BatchSize: &batchSize,
	}
	cur, err := dbService.collectionRefUploadedFiles(instanceID).Find(ctx, filter, &opts)

	if err != nil {
		return fileInfoList, err
	}
	defer cur.Close(ctx)

	fileInfoList = []types.FileInfo{}
	for cur.Next(ctx) {
		var result types.FileInfo
		err := cur.Decode(&result)

		if err != nil {
			return fileInfoList, err
		}

		fileInfoList = append(fileInfoList, result)
	}
	if err := cur.Err(); err != nil {
		return fileInfoList, err
	}

	return fileInfoList, nil
}
