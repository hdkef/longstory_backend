package driver

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURI = os.Getenv("MONGOURI")

func init() {
	godotenv.Load()
}

func InitDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic("DB connection error")
	}
	//TOBEIMPLEMENT
	//ping db conn
	return client
}
