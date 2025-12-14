package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func main() {
	ctx := context.Background()

	// Azurite connection details

	client, err := azblob.NewClientFromConnectionString(
		"DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;QueueEndpoint=http://127.0.0.1:10001/devstoreaccount1;",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.CreateContainer(ctx, "audit", nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.UploadBuffer(context.TODO(), "audit", "2025/12/09/audit.json", []byte("test-audit-log-content"), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Uploaded audit log to Azurite.")

	// Download
	downloadResp, err := client.DownloadStream(ctx, "audit", "2025/12/09/audit.json", nil)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(downloadResp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n--- Retrieved Audit Log ---")
	fmt.Println(string(data))
	fmt.Println("----------------------------")
}
