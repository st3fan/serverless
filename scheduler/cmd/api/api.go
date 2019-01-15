package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	lmbd "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

type TaskRequest struct {
	FunctionName string `json:"functionName"`
	Payload      string `json:"payload"`
	Delay        int64  `json:"delay"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var r TaskRequest
	if err := json.Unmarshal([]byte(request.Body), &r); err != nil {
		return events.APIGatewayProxyResponse{}, errors.Wrap(err, "failed to decode request")
	}

	ses, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.Wrap(err, "failed to get session")
	}

	svc := sqs.New(ses)

	url, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(os.Getenv("TASKS_QUEUE_NAME"))})
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.Wrapf(err, "failed to get queue <%s>", os.Getenv("TASKS_QUEUE_NAME"))
	}

	encodedRequest, err := json.Marshal(r)
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.Wrap(err, "failed to encode request")
	}

	sendMessageInput := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(r.Delay),
		MessageBody:  aws.String(string(encodedRequest)),
		QueueUrl:     url.QueueUrl,
	}

	if _, err := svc.SendMessage(sendMessageInput); err != nil {
		return events.APIGatewayProxyResponse{}, errors.Wrap(err, "failed to send message")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lmbd.Start(Handler)
}
