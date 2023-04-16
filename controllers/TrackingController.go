package controllers

import (
	"TrackingService/utils"
	"cloud.google.com/go/pubsub"
	"context"
	"log"
	"net/http"
	"sync/atomic"
)

const (
	projectID = "trackingservice-383922"
	subID     = "TrackingService"
	topic     = "swipe-track-record"
)

var response *utils.Response

type TrackingControllerInterface interface {
	GetClient()
	SetClient(client *pubsub.Client)
	Ctx()
	SetCtx(ctx context.Context)
}

type TrackingController struct {
	client *pubsub.Client
	ctx    context.Context
}

func (tc *TrackingController) Ctx() context.Context {
	return tc.ctx
}

func (tc *TrackingController) SetCtx(ctx context.Context) {
	tc.ctx = ctx
}

func (tc *TrackingController) GetClient() *pubsub.Client {
	return tc.client
}

func (tc *TrackingController) SetClient(client *pubsub.Client) {
	tc.client = client
}

func (tc *TrackingController) PublishMessage(w http.ResponseWriter, r *http.Request) {

}

func (tc *TrackingController) GetMessagesFromTopic(w http.ResponseWriter, r *http.Request) {
	sub := tc.GetClient().Subscription(subID)

	var received int32
	err := sub.Receive(tc.ctx, func(_ context.Context, msg *pubsub.Message) {
		log.Println(w, "Got message: %q\n", string(msg.Data))
		atomic.AddInt32(&received, 1)
		msg.Ack()
	})
	if err != nil {
		log.Fatalf("sub.Receive:", err)
		response.Format(w, r, true, 418, err)
		return
	}
	log.Println(w, "Received %d messages\n", received)
	response.Format(w, r, false, 200, nil)
	return
}
