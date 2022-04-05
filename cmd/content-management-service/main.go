package main

import (
	"github.com/coneno/logger"
	"github.com/tekenradar/content-service/pkg/dbs/contentdb"
)

func main() {

	conf := InitConfig()

	logger.SetLevel(conf.LogLevel)

	contentDBService := contentdb.NewContentDBService(conf.ContentDBConfig)

	logger.Debug.Println(contentDBService)

}
