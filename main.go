package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Name string
}

func HandleRequest(ctx context.Context, event Event) (string, error) {
	return fmt.Sprintf("Hello %v, I´m a Lambda", event.Name), nil
}

func main() {
	lambda.Start(HandleRequest)
}
