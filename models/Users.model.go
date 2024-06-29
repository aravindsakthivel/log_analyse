package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	Name          string               `bson:"name"`
	UserEmail     string               `bson:"email,unique,required"`
	FilesUploaded []primitive.ObjectID `bson:"filesUploaded"`
}
