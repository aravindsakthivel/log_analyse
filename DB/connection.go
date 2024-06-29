package DB

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	envPkg "github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DBClient *mongo.Client = nil // Database client
var DBWriteClient *mongo.Client = nil

type SDB struct{}

func (db *SDB) Init() error {
	envPkg.Load(".env")

	var mongoUrl string = os.Getenv("DB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))

	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
		return err
	}

	DBClient = client

	pong := client.Ping(ctx, readpref.Primary())

	if pong != nil {
		log.Fatal("Error pinging MongoDB: ", err)
		return pong
	}

	log.Println("Connected to MongoDB")
	return nil
}

func (db *SDB) Health() bool {
	if DBClient == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := DBClient.Ping(ctx, readpref.Primary())
	return err == nil
}

func (db *SDB) Close() {
	if DBClient != nil {
		DBClient.Disconnect(context.Background())
	}
	log.Fatal("Database connection not initialized")
}

func ConnectCL(collection string) (*mongo.Collection, error) {
	envPkg.Load(".env")
	DB_NAME := os.Getenv("DB_NAME")
	if DB_NAME == "" {
		return nil, errors.New("DB_NAME not set")
	}
	log.Println("Connecting to collection: ", collection, " in database: ", DB_NAME)
	if DBClient == nil {
		log.Fatal("Database connection not initialized")
	}

	return DBClient.Database(DB_NAME).Collection(collection), nil
}

func DBClientHealth() bool {
	if DBClient == nil {
		return false
	}
	err := DBClient.Ping(context.Background(), readpref.Primary())
	return err == nil
}

func GetDBClient() *mongo.Client {
	return DBClient
}

func Close() {
	if DBClient != nil {
		DBClient.Disconnect(context.Background())
	}

	log.Fatal("Database connection not initialized")
}
