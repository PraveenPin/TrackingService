package main

import (
	"TrackingService/gclient"
	"TrackingService/routes"
)

func main() {
	app := &gclient.App{}
	ctx := app.GetAppContext()
	client := app.GetPubSubClient(ctx)

	dispatcher := routes.Dispatcher{}
	dispatcher.Init(ctx, client)

	app.CloseClient(client)
}
