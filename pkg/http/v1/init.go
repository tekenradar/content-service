package v1

import "github.com/tekenradar/content-service/pkg/dbs/contentdb"

type HttpEndpoints struct {
	contentDB *contentdb.ContentDBService
	apiKeys   struct {
		readOnly  []string
		readWrite []string
	}
}

func NewHTTPHandler(
	contentDB *contentdb.ContentDBService,
	readOnlyAPIKeys []string,
	readWriteAPIKeys []string,
) *HttpEndpoints {
	return &HttpEndpoints{
		contentDB: contentDB,
		apiKeys: struct {
			readOnly  []string
			readWrite []string
		}{
			readOnly:  readOnlyAPIKeys,
			readWrite: readWriteAPIKeys,
		},
	}
}
