package v1

import "github.com/tekenradar/content-service/pkg/dbs/contentdb"

type HttpEndpoints struct {
	contentDB *contentdb.ContentDBService
}

func NewHTTPHandler(
	contentDB *contentdb.ContentDBService,
) *HttpEndpoints {
	return &HttpEndpoints{
		contentDB: contentDB,
	}
}
