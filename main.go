package main

import (
	"github.com/PraveenPin/TrackingService/clients"
)

func main() {
	app := clients.App{}
	ctx := app.GetAppContext()
	client := app.GetPubSubClient(ctx)
	dynamodbSVC := app.GetDynamoDatabaseClient(app.StartAWSSession())

	dispatcher := Dispatcher{}
	dispatcher.Init(ctx, client, dynamodbSVC)

	app.CloseClient(client)
	return
}
