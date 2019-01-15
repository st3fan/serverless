package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	lmbd "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

// TODO This is copypasted three times. Also the name stupid and inconsistent.
type SubmitRequest struct {
	FunctionName string `json:"functionName"`
	Payload      string `json:"payload"`
	Delay        int64  `json:"delay"`
}

func Handler(ctx context.Context, event events.SQSEvent) error {
	ses, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	if err != nil {
		return errors.Wrap(err, "failed to get session")
	}

	svc := sqs.New(ses)

	url, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(os.Getenv("TASKS_QUEUE_NAME"))})
	if err != nil {
		return errors.Wrapf(err, "failed to get queue <%s>", os.Getenv("TASKS_QUEUE_NAME"))
	}

	client := lambda.New(ses, &aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	for _, message := range event.Records {
		deleteMessageInput := &sqs.DeleteMessageInput{
			QueueUrl:      url.QueueUrl,
			ReceiptHandle: aws.String(message.ReceiptHandle),
		}

		if _, err := svc.DeleteMessage(deleteMessageInput); err != nil {
			log.Println("failed to delete message:", err)
			continue
		}

		var submitRequest SubmitRequest
		if err := json.Unmarshal([]byte(message.Body), &submitRequest); err != nil {
			log.Println("failed to unmarshal request:", err)
			continue
		}

		invokeInput := &lambda.InvokeInput{
			InvocationType: aws.String(lambda.InvocationTypeEvent),
			FunctionName:   aws.String(submitRequest.FunctionName),
			Payload:        []byte(submitRequest.Payload),
		}

		if _, err := client.Invoke(invokeInput); err != nil {
			log.Println("failed to execute lambda:", err)
			continue
		}
	}

	return nil
}

func main() {
	lmbd.Start(Handler)
}
