package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	v1lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"lambda-deploy/packages/slack"
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

func updateFunctionCode(bucketName, functionName, objectName string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		return fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := lambda.NewFromConfig(cfg)

	_, err = client.UpdateFunctionCode(context.TODO(), &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(functionName),
		S3Bucket:     aws.String(bucketName),
		S3Key:        aws.String(objectName),
	})

	if err != nil {
		return fmt.Errorf("failed to update function code, %v", err)
	}

	return nil
}

func sendSlackNotification(functionName string) error {
	fields := []slack.Field{
		{
			Title: "Lambda Function Deploy Complete",
			Value: fmt.Sprintf("FunctionName: %s", functionName),
		},
	}

	attachment := slack.Attachment{Color: "#2eb67d", Fields: fields}
	data := slack.MessageData{Attachments: []slack.Attachment{attachment}}
	log.Printf("Sending Slack message with data: %v\n", data)

	err := slack.SendMessage(data)
	if err != nil {
		log.Printf("Error sending Slack message: %v\n", err)
		return err
	}

	log.Printf("Slack message sent successfully\n")
	return nil
}

func extractS3EventDetails(event S3PutEvent) (string, string, string) {
	bucketName := event.Records[0].S3.Bucket.Name
	functionName := strings.Split(event.Records[0].S3.Object.Key, "/")[0]
	objectName := event.Records[0].S3.Object.Key

	log.Printf("Extracted bucketName: %s", bucketName)
	log.Printf("Extracted functionName: %s", functionName)
	log.Printf("Extracted objectName: %s", objectName)

	return bucketName, functionName, objectName
}

func handler(ctx context.Context, event S3PutEvent) (string, error) {
	bucketName, functionName, objectName := extractS3EventDetails(event)
	log.Printf("function: %s bucket: %s object: %s", functionName, bucketName, objectName)

	err := updateFunctionCode(bucketName, functionName, objectName)
	if err != nil {
		log.Printf("error updating function code: %v", err)
		return "", err
	}

	log.Printf("function code updated successfully")

	err = sendSlackNotification(functionName)
	if err != nil {
		log.Printf("error sending slack notification: %v", err)
		return "", err
	}

	log.Printf("slack notification sent successfully")
	return "Success", nil
}

func main() {
	v1lambda.Start(handler)
}
