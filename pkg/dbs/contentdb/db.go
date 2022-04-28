package contentdb

import (
	"context"
	"time"

	"github.com/coneno/logger"
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
		logger.Error.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Timeout)*time.Second)
	defer cancel()

	err = dbClient.Connect(ctx)
	if err != nil {
		logger.Error.Fatal(err)
	}

	ctx, conCancel := context.WithTimeout(context.Background(), time.Duration(configs.Timeout)*time.Second)
	err = dbClient.Ping(ctx, nil)
	defer conCancel()
	if err != nil {
		logger.Error.Fatal("fail to connect to DB: " + err.Error())
	}

	return &ContentDBService{
		DBClient:     dbClient,
		timeout:      configs.Timeout,
		DBNamePrefix: configs.DBNamePrefix,
	}
}

//new Collection
func (dbService *ContentDBService) collectionRefTickBiteMapInfos(instanceID string) *mongo.Collection {
	return dbService.DBClient.Database(dbService.DBNamePrefix + instanceID + "_contentDB").Collection("tick-bite-map-infos")
}

// DB utils
func (dbService *ContentDBService) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(dbService.timeout)*time.Second)
}
