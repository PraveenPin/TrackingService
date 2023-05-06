package main

import (
	"github.com/PraveenPin/TrackingService/clients"
	"github.com/PraveenPin/TrackingService/routes"
)

func main() {
	app := clients.App{}
	ctx := app.GetAppContext()
	client := app.GetPubSubClient(ctx)
	dynamodbSVC := app.GetDynamoDatabaseClient(app.StartAWSSession())

	dispatcher := routes.Dispatcher{}
	dispatcher.Init(ctx, client, dynamodbSVC)

	app.CloseClient(client)
	return
}
