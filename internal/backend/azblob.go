package backend

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// fetchAzblob downloads a blob from Azure Blob Storage. The storage account
// comes from the URI (azblob://account@container/key) or AZURE_STORAGE_ACCOUNT;
// credentials come from the default Azure credential chain.
func fetchAzblob(ctx context.Context, s Source) ([]byte, error) {
	account := s.Account
	if account == "" {
		account = os.Getenv("AZURE_STORAGE_ACCOUNT")
	}
	if account == "" {
		return nil, fmt.Errorf("azblob: storage account not set (use azblob://account@container/key or AZURE_STORAGE_ACCOUNT)")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("azblob: credential: %w", err)
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", account)
	client, err := azblob.NewClient(serviceURL, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("azblob: client: %w", err)
	}

	resp, err := client.DownloadStream(ctx, s.Container, s.Key, nil)
	if err != nil {
		return nil, fmt.Errorf("azblob: download %s/%s: %w", s.Container, s.Key, err)
	}
	body := resp.Body
	defer body.Close()
	return io.ReadAll(body)
}
