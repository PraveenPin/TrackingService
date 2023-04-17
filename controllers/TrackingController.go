package controllers

import (
	"TrackingService/models"
	"TrackingService/utils"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
)

const (
	projectID = "trackingservice-383922"
	subID     = "TrackingSubscription"
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

	message := models.SwipeRecordMessage{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)
	if err != nil {
		response.Format(w, r, true, 417, err)
		return
	}

	topicClient := tc.client.Topic(topic)
	result := topicClient.Publish(tc.ctx, &pubsub.Message{
		Data: []byte(message.ConvertToString()),
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(tc.ctx)
	if err != nil {
		log.Fatalf("pubsub: result.Get: %v", err)
		response.Format(w, r, true, 418, err)
	}
	log.Println(w, "Published a message; msg ID: \n", id)
	response.Format(w, r, false, 201, nil)

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
