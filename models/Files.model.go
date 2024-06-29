package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Files struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FileName   string             `bson:"fileName"`
	FilePath   string             `bson:"filePath"`
	Date       int64              `bson:"date"`
	InitSearch []string           `bson:"initSearch"`
	Compressed bool               `bson:"compressed"`
	UploadedBy string             `bson:"uploadedBy"`
}
