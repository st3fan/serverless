package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context) (string, error) {
	log.Println("This is hello.Handler")
	return fmt.Sprintf("Hello, %s", os.Getenv("HELLO_NAME")), nil
}

func main() {
	lambda.Start(Handler)
}
