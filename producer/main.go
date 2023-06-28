package main

import (
	"context"
	"fmt"
	"learning-pulsar/domain"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func main() {
	ctx := context.Background()
	err := runSimple(ctx)
	if err != nil {
		panic(err)
	}
}

func runSimple(ctx context.Context) error {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://host.docker.internal:6650",
	})
	if err != nil {
		return err
	}
	defer client.Close()

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic:               "orders",
		BatchingMaxMessages: 10,
		DisableBatching:     false,
	})
	if err != nil {
		return err
	}
	defer producer.Close()

	start := time.Now()

	// 1万笔约24秒
	for i := 0; i < 10000; i++ {
		eventID := uuid.NewString()
		data := domain.CreatedOrderEvent{
			ID:            eventID,
			Sequence:      uint64(i),
			ClientOrderID: eventID,
			Market:        "BTC_USDT",
			Type:          domain.SpotOrderType_Market,
			Price:         decimal.NewFromInt(100),
			Size:          decimal.NewFromInt(100),
			Side:          domain.Side_Buy,
			TimeInForce:   domain.TIF_GTC,
			TakerFeeRate:  decimal.NewFromInt(0),
			MakerFeeRate:  decimal.NewFromInt(0),
			UserID:        123456,
			Source:        "api",
			CreatedAt:     time.Now().UTC(),
		}

		cloudEvent := cloudevents.NewEvent()
		cloudEvent.SetID(eventID)
		cloudEvent.SetSource("trade")
		cloudEvent.SetType("order.created")
		cloudEvent.SetTime(data.CreatedAt)
		err = cloudEvent.SetData(cloudevents.ApplicationJSON, data)
		if err != nil {
			return err
		}

		b, err := cloudEvent.MarshalJSON()
		if err != nil {
			return err
		}

		// _, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		// 	Payload: b,
		// })
		// if err != nil {
		// 	return err
		// }

		producer.SendAsync(ctx, &pulsar.ProducerMessage{
			Payload: b,
		}, func(mi pulsar.MessageID, pm *pulsar.ProducerMessage, err error) {

		})
	}

	fmt.Println("duration:", time.Since(start))
	return nil
}
