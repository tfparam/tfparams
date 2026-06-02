package backend

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// fetchGCS downloads an object from Google Cloud Storage using Application
// Default Credentials.
func fetchGCS(ctx context.Context, s Source) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("gcs: new client: %w", err)
	}
	defer func() { _ = client.Close() }()

	r, err := client.Bucket(s.Bucket).Object(s.Key).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("gcs: read gs://%s/%s: %w", s.Bucket, s.Key, err)
	}
	defer func() { _ = r.Close() }()
	return io.ReadAll(r)
}
