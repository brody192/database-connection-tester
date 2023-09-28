package databases

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Mongo(monoURL string) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	sT := time.Now()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(monoURL))
	if err != nil {
		return time.Since(sT), fmt.Errorf("mongo.Connect error: %w", err)
	}

	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		return time.Since(sT), fmt.Errorf("client.Ping error: %w", err)
	}

	return time.Since(sT), nil
}
