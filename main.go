package main

import (
	"github.com/PraveenPin/TrackingService/clients"
)

const (
	projectID = "trackingservice-383922"
	subID     = "TrackingSubcription"
	topic     = "swipe-track-record"
)

func main() {
	app := clients.App{}
	ctx := app.GetAppContext()
	client := app.GetPubSubClient(ctx)
	dynamodbSVC := app.GetDynamoDatabaseClient(app.StartAWSSession())

	redisClient := app.GetRedisClient()
	dispatcher := Dispatcher{}
	dispatcher.Init(ctx, client, redisClient, dynamodbSVC)

	app.CloseClient(client)
	return
}
