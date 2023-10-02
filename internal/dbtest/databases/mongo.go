package databases

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

func Mongo(monoURL string) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)

	defer cancel()

	sT := time.Now()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(monoURL))
	if err != nil {
		return time.Since(sT), fmt.Errorf("mongo.Connect error: %w", err)
	}

	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		if sse, ok := err.(topology.ServerSelectionError); ok && len(sse.Desc.Servers) > 0 {
			return time.Since(sT), fmt.Errorf("client.Ping error: %w", sse.Desc.Servers[0].LastError)
		}

		return time.Since(sT), fmt.Errorf("client.Ping error: %w", err)
	}

	return time.Since(sT), nil
}
