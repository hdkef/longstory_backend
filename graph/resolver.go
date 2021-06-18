package graph

import "go.mongodb.org/mongo-driver/mongo"

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *mongo.Client
}
