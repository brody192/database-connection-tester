package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/edgedb/edgedb-go"
)

func EdgeDB(edgeDbURL string) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	sT := time.Now()

	client, err := edgedb.CreateClientDSN(ctx, edgeDbURL, edgedb.Options{
		TLSOptions: edgedb.TLSOptions{
			SecurityMode: edgedb.TLSModeInsecure,
		},
	})
	if err != nil {
		return time.Since(sT), err
	}

	defer client.Close()

	if err := client.EnsureConnected(ctx); err != nil {
		return time.Since(sT), fmt.Errorf("client.EnsureConnected() error: %w", err)
	}

	return time.Since(sT), nil
}
