package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/PraveenPin/TrackingService/controllers"
	"github.com/PraveenPin/TrackingService/services"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const PORT = ":8082"

type Dispatcher struct {
}

func HomeEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world :)")
}

func (r *Dispatcher) Init(ctx context.Context, client *pubsub.Client, redisClient *redis.Client, db *dynamodb.DynamoDB) {
	log.Println("Initialize the router")
	router := mux.NewRouter()

	trackingService := services.NewTrackingService(client, ctx, redisClient, db)
	trackingController := controllers.NewTrackingController(trackingService)

	router.StrictSlash(true)
	router.HandleFunc("/", HomeEndpoint).Methods("GET")
	router.HandleFunc("/publishMessage", trackingController.PublishMessage).Methods("POST")

	//router.HandleFunc("/processNewMessages", trackingController.GetMessagesFromTopic).Methods("GET")

	// bind the routes
	http.Handle("/", router)

	log.Println("Add the listener to port", PORT)

	http.ListenAndServe(PORT, nil)
}

func profile(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("test"))
}
