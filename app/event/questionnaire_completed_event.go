package event

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"strconv"
)

const (
	QuestionnaireCompleted = "QUESTIONNAIRE_COMPLETED"
)

type QuestionnaireCompletedEvent struct {
	Name                 string // defines the type of event
	Id                   string
	UserId               string
	StudyId              string
	QuestionnaireId      string
	CompletedAt          string
	RemainingCompletions int
}

func (q *QuestionnaireCompletedEvent) FunctionName() string {
	return q.Name
}

func (q *QuestionnaireCompletedEvent) ToSQSMessage() map[string]*sqs.MessageAttributeValue {
	return map[string]*sqs.MessageAttributeValue{
		"Id": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(q.Id),
		},
		"UserId": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(q.UserId),
		},
		"StudyId": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(q.StudyId),
		},
		"QuestionnaireId": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(q.QuestionnaireId),
		},
		"CompletedAt": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(q.CompletedAt),
		},
		"RemainingCompletions": &sqs.MessageAttributeValue{
			DataType:    aws.String("Number"),
			StringValue: aws.String(strconv.Itoa(q.RemainingCompletions)),
		},
	}
}

/*
	TODO: Would-be handling of a specific QuestionnaireCompleteEvent if that's something that can be invoked by lambda
*/
func (q *QuestionnaireCompletedEvent) HandleEvent(ctx context.Context) error {
	return nil
}
