package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jamesineda/reschedular/app/event"
	"log"
)

const (
	ERROR = "ERROR"
)

func LambdaHandler(ctx context.Context, e event.IncomingEvent) (string, error) {
	//log.Print("Running ƛ %s", ctx.FunctionName)
	if err := event.HandleEvent(ctx, e); err != nil {
		log.Fatalf("received unhandled incoming event: %s", e.Name())
		return ERROR, err
	}

	return fmt.Sprintf("Hello %s!", e.Name()), nil
}

func main() {
	lambda.Start(LambdaHandler)
}
