package event

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"sync"
	"time"
)

// Queue A queue for Asynchronous processing of SQS
type Queue interface {
	Pop() IncomingEvent
	Push(instruction IncomingEvent)
}

type IncomingEvent interface {
	FunctionName() string
	ToSQSMessage() map[string]*sqs.MessageAttributeValue
	HandleEvent(ctx context.Context) error
}
type IncomingEvents []IncomingEvent

// StartAsynchronousEventProcessor a background process that pops events off a queue and sends them out to SQS
func StartAsynchronousEventProcessor(ctx context.Context, wg *sync.WaitGroup, c <-chan bool) {
	wg.Add(1)
	go func() {
		for {
			select {
			case <-c:
				log.Println("SQS queue shutting down")
				wg.Done()
				return

			default:
				eventsQueue := ctx.Value("eventsQueue").(Queue)
				svc := ctx.Value("svc").(*sqs.SQS)
				svcQueueUrl := ctx.Value("scsQueueUrl").(*string)

				if queuedEvent := eventsQueue.Pop(); queuedEvent != nil {
					_, err := svc.SendMessage(&sqs.SendMessageInput{
						DelaySeconds:      aws.Int64(10),
						MessageAttributes: queuedEvent.ToSQSMessage(),
						MessageBody:       aws.String("Information about current NY Times fiction bestseller for week of 12/11/2016."),
						QueueUrl:          svcQueueUrl,
					})

					// If we fail to submit the event to SQS, log the error. A better solution would be some sort of automatic
					// retry and after N attempts, store as a failure in an "events" table or something. This would allow us to
					// manually trigger a message via a console or something after the issue has been resolved.
					if err != nil {
						log.Printf("failed to submit event %s to SQS: %s", queuedEvent.FunctionName(), err)
					}
				}

				time.Sleep(10 * time.Millisecond) // reduce CPU usage, less spam
			}
		}
	}()
}
