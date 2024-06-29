package controller

import (
	"context"
	"log"
	"log_analyse/DB"
	"log_analyse/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SFilesInfo struct {
	ID         primitive.ObjectID
	FileName   string
	FilePath   string
	Date       string
	InitSearch []string
	Compressed bool
}

type IFilesCL interface {
	SetCollection() error
	SetFiles(filePath string, fileName string, email string) (primitive.ObjectID, error)
	GetUnCompressedFile() (bool, []string, error)
}

type SFilesCL struct {
	Collection *mongo.Collection
}

func (f *SFilesCL) setCollection() error {

	if f.Collection != nil {
		return nil
	}

	cl, err := DB.ConnectCL("files")

	if err != nil {
		log.Print("Error connecting to files collection: ", err)
		return err
	}

	f.Collection = cl

	return nil
}

func (f *SFilesCL) SetFiles(filePath string, fileName string, email string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var fileId primitive.ObjectID = primitive.NewObjectID()
	var timeStamp = time.Now().UTC().UnixMilli()
	var fileModel models.Files = models.Files{
		ID:         fileId,
		FileName:   fileName,
		FilePath:   filePath,
		Date:       timeStamp,
		InitSearch: []string{},
		UploadedBy: email,
		Compressed: false,
	}

	_, err := f.Collection.InsertOne(ctx, fileModel)

	if err != nil {
		log.Printf("Error inserting file %s %s", fileName, err)
		return fileId, err
	}

	log.Printf("Inserted file info %s %s", fileName, fileId)

	return fileId, nil
}

func (f *SFilesCL) GetUnCompressedFile() (bool, []string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	filter := bson.D{{Key: "compressed", Value: false}}

	projection := bson.D{{Key: "filePath", Value: 1}}

	cursor, err := f.Collection.Find(ctx, filter, options.Find().SetProjection(projection).SetLimit(3))

	if err != nil {
		log.Printf("Error getting unCompressed files %s ", err)
		return false, []string{}, err
	}

	var results []struct {
		FilePath string `bson:"filePath"`
	}

	crErr := cursor.All(ctx, &results)

	if crErr != nil {
		log.Printf("Error decoding unCompressed files %s ", err)
		return false, []string{}, err
	}

	log.Print("Response from DB ", results)

	var filePaths []string

	for _, result := range results {
		filePaths = append(filePaths, result.FilePath)
	}

	if len(filePaths) == 0 {
		return false, []string{}, nil
	}

	return true, filePaths, nil
}
