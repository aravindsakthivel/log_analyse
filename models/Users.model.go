package models

type Users struct {
	Name          string   `bson:"name"`
	UserEmail     string   `bson:"email,unique,required"`
	FilesUploaded []string `bson:"filesUploaded"`
}
