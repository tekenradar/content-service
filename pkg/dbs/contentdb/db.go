package contentdb

import (
	"context"
	"log"
	"time"

	"github.com/tekenradar/content-service/pkg/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContentDBService struct {
	DBClient     *mongo.Client
	timeout      int
	DBNamePrefix string
}

func NewContentDBService(configs types.DBConfig) *ContentDBService {
	var err error
	dbClient, err := mongo.NewClient(
		options.Client().ApplyURI(configs.URI),
		options.Client().SetMaxConnIdleTime(time.Duration(configs.IdleConnTimeout)*time.Second),
		options.Client().SetMaxPoolSize(configs.MaxPoolSize),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Timeout)*time.Second)
	defer cancel()

	err = dbClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ctx, conCancel := context.WithTimeout(context.Background(), time.Duration(configs.Timeout)*time.Second)
	err = dbClient.Ping(ctx, nil)
	defer conCancel()
	if err != nil {
		log.Fatal("fail to connect to DB: " + err.Error())
	}

	return &ContentDBService{
		DBClient:     dbClient,
		timeout:      configs.Timeout,
		DBNamePrefix: configs.DBNamePrefix,
	}
}

//new Collection
func (dbService *ContentDBService) collectionRefStudyInfos(instanceID string) *mongo.Collection {
	return dbService.DBClient.Database(dbService.DBNamePrefix + instanceID + "_contentDB").Collection("map-infos")
}

// DB utils
//func (dbService *StudyDBService) getContext() (ctx context.Context, cancel context.CancelFunc) {
//	return context.WithTimeout(context.Background(), time.Duration(dbService.timeout)*time.Second)
//}
