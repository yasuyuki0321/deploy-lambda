package main

import (
	"context"
	"log"
	"strings"

	v1Lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type S3PutEvent struct {
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func handler(ctx context.Context, event S3PutEvent) error {
	bucketName := event.Records[0].S3.Bucket.Name
	functionName := strings.Split(event.Records[0].S3.Object.Key, "/")[0]
	objectName := strings.Split(event.Records[0].S3.Object.Key, "/")[1]

	log.Printf("bucket: %s, function: %s object: %s", bucketName, functionName, objectName)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := lambda.NewFromConfig(cfg)

	_, err = client.UpdateFunctionCode(context.TODO(), &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(functionName),
		S3Bucket:     aws.String(bucketName),
		S3Key:        aws.String(objectName),
	})

	if err != nil {
		log.Fatalf("unable to update function code, %v", err)
	}

	return nil
}

func main() {
	v1Lambda.Start(handler)
}
