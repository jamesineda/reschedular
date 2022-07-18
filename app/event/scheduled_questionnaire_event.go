package event

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"time"
)

const (
	ScheduledQuestionnaire = "SCHEDULED_QUESTIONNAIRE"
	Pending                = "pending"
	Completed              = "completed"
)

var ErrMaxAttemptsReached = fmt.Errorf("maximum number of results reached for questionnaire ")
var ErrScheduledQuestionnaireIsAlreadyCompleted = fmt.Errorf("scheduled questionnaire marked as completed")

// ScheduledQuestionnaireEvent provides an interface to handle ScheduledQuestionnaire events via SQS message transmission
// Lambda call handling can be added by implementing the HandleEvent method
type ScheduledQuestionnaireEvent struct {
	Name            string // defines the type of event
	Id              string
	ParticipantId   string
	QuestionnaireId string
	Status          string
	ScheduledAt     time.Time
}

func (q *ScheduledQuestionnaireEvent) FunctionName() string {
	return q.Name
}

func (q *ScheduledQuestionnaireEvent) ToSQSMessage() map[string]*sqs.MessageAttributeValue {
	return map[string]*sqs.MessageAttributeValue{
		"Id": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(q.Id),
		},
		"ParticipantId": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(q.ParticipantId),
		},
		"QuestionnaireId": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(q.QuestionnaireId),
		},
		"Status": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(q.Status),
		},
		"ScheduledAt": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(q.ScheduledAt.Format(time.RFC3339)),
		},
	}
}

// HandleEvent No specific handling for this function from a Lambda call just yet
func (event *ScheduledQuestionnaireEvent) HandleEvent(ctx context.Context) (err error) {
	return
}
