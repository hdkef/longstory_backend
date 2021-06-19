package main

import (
	"log"
	"longstory/graph"
	"longstory/graph/generated"
	"longstory/rest"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

const defaultPort = "8080"

func init() {
	godotenv.Load()
}

func main() {

	// db := driver.InitDB()
	db := &mongo.Client{}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	router.HandleFunc("/video", rest.UploadVideo(db))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
