package controllers

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/PraveenPin/SwipeMeter/utils"
	"github.com/PraveenPin/TrackingService/models"
	"github.com/PraveenPin/TrackingService/services"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
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
	RedisClient()
	SetRedisClient(redisClient *redis.Client)
}

type TrackingController struct {
	trackingService *services.TrackingService
}

func NewTrackingController(trackingService *services.TrackingService) *TrackingController {
	return &TrackingController{trackingService: trackingService}
}

func (tc *TrackingController) PublishMessage(w http.ResponseWriter, r *http.Request) {

	message := models.SwipeRecordMessage{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)
	if err != nil {
		response.Format(w, r, true, 417, err)
		return
	}

	_, pub_err := tc.trackingService.PublishMessageService(topic, message)

	if pub_err != nil {
		log.Println("pubsub: result.Get: %v", pub_err)
		response.Format(w, r, true, 418, pub_err)
	}

	response.Format(w, r, false, 201, nil)

}
