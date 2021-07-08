package certificate

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// AzureBlob represents a AzureBlob connection to a named container
type AzureBlob struct {
	accountName   string
	accessKey     string
	containerName string
	containerURL  azblob.ContainerURL
}

// NewAzureBlob creates an instance of AzureBlob
// accountName is the AZURE_STORAGE_ACCOUNT
// accessKey is the AZURE_STORAGE_ACCESS_KEY
// containerName is the container name
func NewAzureBlob(accountName, accessKey, containerName string) (*AzureBlob, error) {
	if len(accountName) == 0 || len(accessKey) == 0 || len(containerName) == 0 {
		return nil, fmt.Errorf("Either the AZURE_STORAGE_ACCOUNT, AZURE_STORAGE_ACCESS_KEY or Container name cannot be empty")
	}

	return &AzureBlob{
		accountName:   accountName,
		accessKey:     accessKey,
		containerName: containerName,
	}, nil
}

func (ab *AzureBlob) connect() error {
	credential, err := azblob.NewSharedKeyCredential(ab.accountName, ab.accessKey)
	if err != nil {
		return fmt.Errorf("Invalid credentials with error: " + err.Error())
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", ab.accountName, ab.containerName),
	)

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	ab.containerURL = azblob.NewContainerURL(*URL, p)

	return nil
}

// Download a file from blob
func (ab *AzureBlob) Download(path, filename string) ([]byte, error) {
	ab.connect()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*20))
	defer cancel()

	blobURL := ab.containerURL.NewBlockBlobURL(fmt.Sprintf("%s/%s", path, filename))
	downloadResponse, err := blobURL.Download(
		ctx,
		0,
		azblob.CountToEnd,
		azblob.BlobAccessConditions{},
		false,
		azblob.ClientProvidedKeyOptions{},
	)
	if err != nil {
		return []byte(""), err
	}

	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})

	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	if err != nil {
		return []byte(""), err
	}

	return downloadedData.Bytes(), nil
}
