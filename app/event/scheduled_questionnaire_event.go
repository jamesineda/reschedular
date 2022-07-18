package event

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jamesineda/reschedular/app/db"
	"github.com/jamesineda/reschedular/app/models"
	"github.com/jamesineda/reschedular/app/utils"
	"log"
	"time"
)

const (
	ScheduledQuestionnaire = "SCHEDULED_QUESTIONNAIRE"
	Pending                = "pending"
	Completed              = "completed"
)

var ErrMaxAttemptsReached = fmt.Errorf("maximum number of results reached for questionnaire ")
var ErrScheduledQuestionnaireIsAlreadyCompleted = fmt.Errorf("scheduled questionnaire marked as completed")

/*
	The readme asks me to specifically implement the function of creating a "scheduled questionnaire", but the code
	provided the skeleton code for a QuestionnaireCompletedEvent. I'm assuming the latter is the message payload I need
	to process asynchronously. I'm also assuming a "scheduled questionnaire" event exists and that this microservice
	receives this via an AWS lambda call as well.
*/
type ScheduledQuestionnaireEvent struct {
	Name            string // defines the type of event
	Id              string
	ParticipantId   string
	QuestionnaireId string
	Status          string

	// this gets set duration the scheduled_questionnaire message
	ScheduledAt time.Time
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

func (event *ScheduledQuestionnaireEvent) HandleEvent(ctx context.Context) (err error) {
	dbConn := ctx.Value("db").(db.Client)
	timer := ctx.Value("timer").(utils.Timer)
	idGenny := ctx.Value("idGenny").(utils.IdGenny)

	defer func() {
		eventsQueue := ctx.Value("eventsQueue").(Queue)
		event.handleDeferFunc(err, eventsQueue, timer, idGenny)
	}()

	//	2. Determine if a new questionnaire schedule should be saved to the database.
	questionnaireRow, err := dbConn.GetById(event.QuestionnaireId, &models.Questionnaire{})
	if err != nil {
		return fmt.Errorf("failed to get Questionnaire (id: %s) from database: %v", event.QuestionnaireId, err)
	}
	questionnaire := questionnaireRow.(*models.Questionnaire)

	participantRow, err := dbConn.GetById(event.ParticipantId, &models.Participant{})
	if err != nil {
		return fmt.Errorf("failed to get participant (id: %s) from database: %v", event.ParticipantId, err)
	}
	participant := participantRow.(*models.Participant)

	// I'm assuming a scheduled_questionnaire of status "completed" at this point denotes some sort of delay in processing
	// and that we should handle that as if the participant has already submitted their results
	if event.Status == Completed {
		return ErrScheduledQuestionnaireIsAlreadyCompleted
	}

	existingResultsArgs := db.Filters{
		{"questionnaire_id", "=", questionnaire.Id},
		{"participant_id", "=", participant.Id},
		{"questionnaire_schedule_id", "=", event.Id}}

	var existingResults models.QuestionnaireResults
	err = dbConn.GetList(&existingResults, existingResultsArgs)
	if err != nil {
		return fmt.Errorf("failed to query existing_results (questionnaire_id: %s, participant_id: %s) from database: %v",
			event.QuestionnaireId, event.ParticipantId, err)
	}

	// checks if we've already received the maximum number of results from the participant for the requested questionnaire
	if !questionnaire.CanAttempt(existingResults.Count()) {
		log.Printf("maximum number of results reached for questionnaire (id: %s, participant_id %s)",
			event.QuestionnaireId, event.ParticipantId)
		return ErrMaxAttemptsReached
	}

	//	3. If so, save one in the database, and push a new message to SQS that a new schedule has been created.
	scheduledQuestionnaire := models.ScheduledQuestionnaire{
		Id:              event.Id,
		QuestionnaireId: event.QuestionnaireId,
		ParticipantId:   event.ParticipantId,
		ScheduledAt:     timer.GetTimeNow().Add(questionnaire.GetHoursBetweenAttemptsDuration()),
		Status:          sql.NullString{Valid: true, String: Pending},
	}

	// set the scheduled time we generate (based off the abstract questionnaire configuration record)
	event.ScheduledAt = scheduledQuestionnaire.ScheduledAt

	// attempt to insert the scheduled_questionnaire into the database
	// I'm going to assume updating a scheduled_questionnaire record would be handled in a separate update event? Presumably
	// but whatever process consumes the QuestionnaireComplete message that this microservices pushes to SQS?
	err = dbConn.Create(&scheduledQuestionnaire)
	return
}

func (event *ScheduledQuestionnaireEvent) handleDeferFunc(err error, eventsQueue Queue, timer utils.Timer, idGenny utils.IdGenny) {
	switch err {
	case nil:
		// pops the scheduled_questionnaire created message onto the events queue for asynchronous SQS transmission
		eventsQueue.Push(event)

	//	4. If not, push a new message to SQS that the user has completed all of their alloted scheduled questionnaires.
	// so, from this, I'm guessing the two scenarios for this would be if it's already completed, or we've reached our
	// maxiumum number of attempts
	case ErrMaxAttemptsReached, ErrScheduledQuestionnaireIsAlreadyCompleted:
		eventsQueue.Push(&QuestionnaireCompletedEvent{
			Id:                   idGenny.GenerateId(),
			UserId:               event.ParticipantId,
			StudyId:              "SUPERSECRET",
			QuestionnaireId:      event.QuestionnaireId,
			CompletedAt:          timer.GetTimeNow().Format(time.RFC3339),
			RemainingCompletions: 0,
		})

	default:
		// unexpected errors handled here, log and cry about it loudly!
		log.Fatalf("failed to process scheduled questionnaire event: %s", err)
	}
}
