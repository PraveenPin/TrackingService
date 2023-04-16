package routes

import (
	"TrackingService/controllers"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Dispatcher struct {
}

func HomeEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world :)")
}

func (r *Dispatcher) Init(ctx context.Context, client *pubsub.Client) {
	log.Println("Initialize the router")
	router := mux.NewRouter()

	trackingController := &controllers.TrackingController{}
	trackingController.SetClient(client)
	trackingController.SetCtx(ctx)

	router.StrictSlash(true)
	router.HandleFunc("/", HomeEndpoint).Methods("GET")
	// User Resource
	//userRoutes := router.PathPrefix("/users").Subrouter()
	router.HandleFunc("/processNewMessages", trackingController.GetMessagesFromTopic).Methods("GET")
	router.HandleFunc("/publishMessage", trackingController.PublishMessage).Methods("POST")

	//Authenticate
	//userRoutes.HandleFunc("/authenticate", UserController.Authenticate).Methods("POST")

	// bind the routes
	http.Handle("/", router)

	log.Println("Add the listener to port 8080")

	//serve
	http.ListenAndServe(":8080", nil)
}

func profile(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("test"))
}
