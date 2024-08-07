package contentdb

import (
	"time"

	"github.com/tekenradar/content-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (dbService *ContentDBService) CreateIndexLPPInfos(instanceID string) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	_, err := dbService.collectionLPP(instanceID).Indexes().CreateOne(
		ctx, mongo.IndexModel{
			Keys: bson.M{
				"pid": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	)
	return err
}

func (dbService *ContentDBService) GetLPPParticipant(instanceID string, pid string) (lppParticipant types.LPPParticipant, err error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{"pid": pid}

	elem := types.LPPParticipant{}
	err = dbService.collectionLPP(instanceID).FindOne(ctx, filter).Decode(&elem)
	return elem, err
}

func (dbService *ContentDBService) AddLPPParticipant(instanceID string, lppParticipant types.LPPParticipant) (string, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	res, err := dbService.collectionLPP(instanceID).InsertOne(ctx, lppParticipant)
	if err != nil {
		return "", err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), err
}

func (dbService *ContentDBService) ReplaceLPPParticipant(instanceID string, lppParticipant types.LPPParticipant) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{"pid": lppParticipant.PID}
	_, err := dbService.collectionLPP(instanceID).ReplaceOne(ctx, filter, lppParticipant)
	if err != nil {
		return err
	}
	return err
}

func (dbService *ContentDBService) FindUninvitedLPPParticipants(instanceID string) ([]types.LPPParticipant, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{"invitationSentAt": bson.M{"$lte": time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)}}
	batchSize := int32(32)
	opts := options.FindOptions{
		BatchSize: &batchSize,
	}
	cur, err := dbService.collectionLPP(instanceID).Find(ctx, filter, &opts)

	if err != nil {
		return []types.LPPParticipant{}, err
	}
	defer cur.Close(ctx)

	participants := []types.LPPParticipant{}
	for cur.Next(ctx) {
		var result types.LPPParticipant
		err := cur.Decode(&result)

		if err != nil {
			return []types.LPPParticipant{}, err
		}

		participants = append(participants, result)
	}
	if err := cur.Err(); err != nil {
		return []types.LPPParticipant{}, err
	}

	return participants, nil
}

func (dbService *ContentDBService) UpdateLPPParticipantInvitationSentAt(instanceID string, pid string, invitationSentAt time.Time) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{"pid": pid}
	update := bson.M{"$set": bson.M{"invitationSentAt": invitationSentAt}}
	_, err := dbService.collectionLPP(instanceID).UpdateOne(ctx, filter, update)
	return err
}

func (dbService *ContentDBService) UpdateLPPParticipantSubmissions(instanceID string, pid string, submissions map[string]time.Time) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{"pid": pid}
	update := bson.M{"$set": bson.M{"submissions": submissions}}
	_, err := dbService.collectionLPP(instanceID).UpdateOne(ctx, filter, update)
	return err
}
