package tester

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
)

// WaitForDeliveryEvent listens for Pub/Sub messages on a subscription and
// returns the first one whose raw payload contains testName.
//
// projectID      – GCP project that owns the Pub/Sub subscription
// subscriptionID – name of the subscription to listen on
// testName       – keyword to look for in the raw message data
// timeout        – overall time to wait before giving up
func WaitForDeliveryEvent(projectID, subscriptionID, testName string, timeout time.Duration) (*pubsub.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub client error: %w", err)
	}
	defer client.Close()

	sub := client.Subscription(subscriptionID)

	var matchedMsg *pubsub.Message
	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		dataStr := string(m.Data)

		//log.Printf("🔔 Received Pub/Sub message ID: %s", m.ID)
		//log.Printf("📦 Raw payload: %q", dataStr)

		if strings.Contains(dataStr, testName) {
			log.Printf("Matched test name %q", testName)
			matchedMsg = m
			m.Ack()
			cancel() // stop Receive loop immediately
			return
		}

		// Not a match – let Pub/Sub redeliver later
		m.Nack()
	})

	// err will be context.Canceled when we call cancel() on a match – ignore it
	if err != nil && err != context.Canceled {
		return nil, fmt.Errorf("receive error: %w", err)
	}

	if matchedMsg == nil {
		return nil, fmt.Errorf("no matching message found for test name %q", testName)
	}

	return matchedMsg, nil
}
