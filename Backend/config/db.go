package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var Client *mongo.Client
func ConnectDB() *mongo.Client {
	uri := "mongodb+srv://expense_user:Expense123@cluster0.0qvxgbx.mongodb.net/?appName=Cluster0"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	Client = client
	fmt.Println("✅ Connected to MongoDB")
	return client
	
}