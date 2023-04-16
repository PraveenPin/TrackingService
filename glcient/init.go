package glcient

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"log"
)

const (
	projectID = "trackingservice-383922"
	subID     = "TrackingService"
	topic     = "swipe-track-record"
)

type App struct {
}

func (a *App) GetAppContext() context.Context {
	fmt.Println("Initialising App Context")
	return context.Background()
}

func (a *App) GetPubSubClient(ctx context.Context) *pubsub.Client {

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("glcient.NewClient:", err)
		return nil
	}
	fmt.Println("Pub/Sub client obtained")
	return client
}

func (a *App) CloseClient(client *pubsub.Client) {
	fmt.Println("Will close the pub/sub client shortly")
	defer client.Close()
}
