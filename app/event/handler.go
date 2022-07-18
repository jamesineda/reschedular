package event

import (
	"context"
	"fmt"
)

type IncomingEvent interface {
	Name() string
}

func HandleEvent(ctx context.Context, event IncomingEvent) error {
	switch event.Name() {
	case QuestionnaireCompleted:
	default:
		return fmt.Errorf("received unhandled event: %s", event.Name())
	}
	return nil
}
