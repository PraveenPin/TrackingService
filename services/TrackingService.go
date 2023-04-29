package services

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/PraveenPin/TrackingService/models"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-redis/redis/v8"
	"log"
)

const (
	subID = "TrackingSubscription"
)

type TrackingService struct {
	client      *pubsub.Client
	ctx         context.Context
	redisClient *redis.Client
	db          *dynamodb.DynamoDB
}

func NewTrackingService(client *pubsub.Client, ctx context.Context, redisClient *redis.Client, db *dynamodb.DynamoDB) *TrackingService {
	return &TrackingService{client: client, ctx: ctx, redisClient: redisClient, db: db}
}

func (ts *TrackingService) PublishMessageService(topic string, message models.SwipeRecordMessage) (bool, error) {
	topicClient := ts.client.Topic(topic)
	result := topicClient.Publish(ts.ctx, &pubsub.Message{
		Data: []byte(message.ConvertToString()),
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ts.ctx)
	if err != nil {
		return false, err
	}
	log.Printf("Published a message; msg ID: \n", id)
	return true, nil
}
