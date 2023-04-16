package TrackingService

import (
	"TrackingService/glcient"
	"TrackingService/routes"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

func pullMsgs(w io.Writer, testPublish bool) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("glcient.NewClient: %v", err)
	}
	//its like finally in other languages, it calls close after exiting the function
	defer client.Close()

	sub := client.Subscription(subID)
	//topic := client.Topic(topic)

	//if testPublish {
	//	// Publish 10 messages on the topic.
	//	var results []*glcient.PublishResult
	//	for i := 0; i < 10; i++ {
	//		res := topic.Publish(ctx, &glcient.Message{
	//			Data: []byte(fmt.Sprintf("hello world #%d", i)),
	//		})
	//		results = append(results, res)
	//	}
	//
	//	// Check that all messages were published.
	//	for _, r := range results {
	//		_, err := r.Get(ctx)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	//defer cancel()

	var received int32
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
		atomic.AddInt32(&received, 1)
		msg.Ack()
		if received == 10 {
			cancel()
		}
	})
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}
	fmt.Fprintf(w, "Received %d messages\n", received)

	return nil
}

func main() {
	app := &glcient.App{}
	ctx := app.GetAppContext()
	client := app.GetPubSubClient(ctx)

	dispatcher := routes.Dispatcher{}
	dispatcher.Init(ctx, client)

	app.CloseClient(client)
}
