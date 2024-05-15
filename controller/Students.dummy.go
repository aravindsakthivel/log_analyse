package controller

import (
	"context"
	"log"
	"log_analyse/DB"
	"log_analyse/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type IStudentCL interface {
}

type SStudentCL struct {
	Collection *mongo.Collection
}

func (s *SStudentCL) SetCollection() error {
	cl, err := DB.ConnectCL("students")

	if err != nil {
		log.Fatal("Error connecting to collection: ", err)
		return err
	}
	s.Collection = cl

	return nil
}

func (s *SStudentCL) CreateStudent() {

	student := models.Student{
		Name:    "John Doe",
		Age:     21,
		Email:   "john.doe@example.com",
		PhoneNo: "123-456-7890",
	}

	insertResult, err := s.Collection.InsertOne(context.TODO(), student)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Inserted student with ID %s\n", insertResult.InsertedID)
}
