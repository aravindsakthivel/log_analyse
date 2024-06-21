package controller

import (
	"context"
	"fmt"
	"log"
	"log_analyse/DB"
	"log_analyse/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SUserInfo struct {
	Name          string
	UserEmail     string
	FilesUploaded []string
}

type IUsersCL interface {
	SetCollection() error
	CreateUser(user SUserInfo) error
	PushFileName(email string, file string) error
}

type SUsersCL struct {
	Collection *mongo.Collection
}

func (s *SUsersCL) setCollection() error {

	if s.Collection != nil {
		return nil
	}

	cl, err := DB.ConnectCL("users")

	if err != nil {
		log.Print("Error connecting to collection: ", err)
		return err
	}
	s.Collection = cl
	return nil
}

func (s *SUsersCL) CreateUser(user SUserInfo) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var userModel models.Users = models.Users{
		Name:          user.Name,
		UserEmail:     user.UserEmail,
		FilesUploaded: user.FilesUploaded,
	}

	var result models.Users

	err := s.Collection.FindOne(ctx, bson.D{{Key: "email", Value: userModel.UserEmail}}).Decode(&result)

	if err != nil && err != mongo.ErrNoDocuments {
		log.Print("Error finding user: ", err)
		return err
	}

	if result.UserEmail == userModel.UserEmail {
		return fmt.Errorf("user with email %s already exists", userModel.UserEmail)
	}

	insertResult, err := s.Collection.InsertOne(ctx, userModel)
	if err != nil {
		log.Printf("Error inserting user: %s %s ", user.UserEmail, err)
		return err
	}

	log.Printf("Inserted student with ID %s\n", insertResult.InsertedID)
	return nil
}

func (s *SUsersCL) PushFileName(email string, file string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var result models.Users

	err := s.Collection.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&result)

	if err != nil && err != mongo.ErrNoDocuments {
		log.Print("Error finding user: ", err)
		return err
	}

	if result.UserEmail == email {
		return fmt.Errorf("user with email %s already exists", email)
	}

	filter := bson.D{{Key: "email", Value: email}}

	update := bson.D{{Key: "$push", Value: bson.D{{
		Key: "filesUploaded", Value: file,
	}}}}

	updateResult, err := s.Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Printf("Error inserting file name %s %s %s ", email, file, err)
		return err
	}

	log.Printf("Inserted file %s %s %d ", email, file, updateResult.ModifiedCount)

	return nil
}
