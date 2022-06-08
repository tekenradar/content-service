package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type FileInfo struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Path       string             `bson:"path,omitempty"`
	UploadedAt int64              `bson:"uploadedAt,omitempty"`
	FileType   string             `bson:"fileType,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Size       int32              `bson:"size,omitempty"`
}
