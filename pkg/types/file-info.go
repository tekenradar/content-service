package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type FileInfo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Path        string             `bson:"path,omitempty" json:"path,omitempty"`
	UploadedAt  int64              `bson:"uploadedAt,omitempty" json:"uploadedAt,omitempty"`
	FileType    string             `bson:"fileType,omitempty" json:"fileType,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Size        int32              `bson:"size,omitempty" json:"size,omitempty"`
	Label       string             `form:"label" bson:"label,omitempty" json:"label,omitempty"`
	Description string             `form:"description" bson:"description,omitempty" json:"description,omitempty"`
}
