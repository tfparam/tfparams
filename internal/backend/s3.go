package backend

import (
	"context"
	"fmt"
	"io"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// fetchS3 downloads an object from S3. Region and credentials come from the AWS
// default credential chain (AWS_REGION, ~/.aws/config, IAM role, ...).
func fetchS3(ctx context.Context, s Source) ([]byte, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("s3: load aws config: %w", err)
	}
	client := s3.NewFromConfig(cfg)
	out, err := client.GetObject(ctx, &s3.GetObjectInput{Bucket: &s.Bucket, Key: &s.Key})
	if err != nil {
		return nil, fmt.Errorf("s3: get s3://%s/%s: %w", s.Bucket, s.Key, err)
	}
	defer out.Body.Close()
	return io.ReadAll(out.Body)
}
