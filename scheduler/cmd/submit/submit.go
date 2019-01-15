package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

type SubmitRequest struct {
	FunctionName string `json:"functionName"`
	Payload      string `json:"payload"`
	Delay        int64  `json:"delay"`
}

func Handler(ctx context.Context, req SubmitRequest) error {
	if len(req.Payload) > 8192 {
		return errors.New("invalid message parameter (too large)")
	}

	if req.Delay <= 0 || req.Delay >= 900 {
		return errors.New("invalid delay parameter (out of range)")
	}

	ses, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	if err != nil {
		return errors.Wrap(err, "failed to get session")
	}

	svc := sqs.New(ses)

	url, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(os.Getenv("TASKS_QUEUE_NAME"))})
	if err != nil {
		return errors.Wrapf(err, "failed to get queue <%s>", os.Getenv("TASKS_QUEUE_NAME"))
	}

	encodedRequest, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "failed to encode request")
	}

	sendMessageInput := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(req.Delay),
		MessageBody:  aws.String(string(encodedRequest)),
		QueueUrl:     url.QueueUrl,
	}

	if _, err := svc.SendMessage(sendMessageInput); err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
