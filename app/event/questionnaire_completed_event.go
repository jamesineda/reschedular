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
	"strconv"
	"time"
)

const (
	QuestionnaireCompleted = "QUESTIONNAIRE_COMPLETED"
)

var ErrAdhocQuestionnaireCompleted = fmt.Errorf("an adhoc questionnaire was completed")

// QuestionnaireCompletedEvent provides an interface to handle QuestionnaireCompleted events via SQS message transmission
// and Lambda call
type QuestionnaireCompletedEvent struct {
	Name                 string // defines the type of event
	Id                   string
	UserId               string
	StudyId              string
	QuestionnaireId      string
	CompletedAt          string
	RemainingCompletions int
}

// GetCompletedAt I've not seen a completedAt time format, so let's pretend all timestamps sent via APIs is in RFC3339
func (q *QuestionnaireCompletedEvent) GetCompletedAt() time.Time {
	t, _ := time.Parse(q.CompletedAt, time.RFC3339)
	return t
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
func (event *QuestionnaireCompletedEvent) HandleEvent(ctx context.Context) (err error) {
	dbConn := ctx.Value("db").(db.Client)
	idGenny := ctx.Value("idGenny").(utils.IdGenny)
	var scheduledQuestionnaire *models.ScheduledQuestionnaire = nil

	defer func() {
		eventsQueue := ctx.Value("eventsQueue").(Queue)
		event.handleDeferFunc(err, eventsQueue, scheduledQuestionnaire)
	}()

	//	2. Determine if a new questionnaire schedule should be saved to the database.
	questionnaireRow, err := dbConn.GetById(event.QuestionnaireId, &models.Questionnaire{})
	if err != nil {
		err = fmt.Errorf("failed to get Questionnaire (id: %s) from database: %v", event.QuestionnaireId, err)
		return
	}
	questionnaire := questionnaireRow.(*models.Questionnaire)

	participantRow, err := dbConn.GetById(event.UserId, &models.Participant{})
	if err != nil {
		// making sure the participant exists in the database
		err = fmt.Errorf("failed to get participant (id: %s) from database: %v", event.UserId, err)
		return
	}
	participant := participantRow.(*models.Participant)

	// checking to see if the questionnaire relates to a scheduled_questionnaire
	scheduledQuestionnairesArgs := db.Filters{
		{"questionnaire_id", "=", questionnaire.Id},
		{"participant_id", "=", participant.Id}}
	var scheduledQuestionnaires models.ScheduledQuestionnaires
	err = dbConn.GetList(&scheduledQuestionnaires, scheduledQuestionnairesArgs)

	switch err {
	case nil:
		existingResultsArgs := db.Filters{
			{"questionnaire_id", "=", questionnaire.Id},
			{"participant_id", "=", participant.Id},
			{"questionnaire_schedule_id", "=", event.Id}}

		var existingResults models.QuestionnaireResults
		err = dbConn.GetList(&existingResults, existingResultsArgs)
		if err != nil && err != sql.ErrNoRows {
			err = fmt.Errorf("failed to query existing_results (questionnaire_id: %s, participant_id: %s) from database: %v",
				event.QuestionnaireId, event.UserId, err)
			return
		}

		// assuming remaining completions is the number of scheduled questionnaires a participant has left to complete, then
		// if zero the service deems questionnaire as complete. I'm also taking into consideration the maximum number of attempts
		//for a particular questionnaire (based off the number of results?)
		if event.RemainingCompletions == 0 || !questionnaire.CanAttempt(existingResults.Count()) {
			log.Printf("maximum number of results reached for questionnaire (id: %s, participant_id %s)",
				event.QuestionnaireId, event.UserId)
			err = ErrMaxAttemptsReached
			return
		}

		//	3. If so, save one in the database, and push a new message to SQS that a new schedule has been created.
		// At this point, we h
		scheduledQuestionnaire = &models.ScheduledQuestionnaire{
			Id:              idGenny.GenerateId(),
			QuestionnaireId: event.QuestionnaireId,
			ParticipantId:   event.UserId,
			ScheduledAt:     event.GetCompletedAt().Add(questionnaire.GetHoursBetweenAttemptsDuration()),
			Status:          sql.NullString{Valid: true, String: Pending},
		}

		// attempt to insert the scheduled_questionnaire into the database
		// I'm going to assume updating a scheduled_questionnaire record would be handled in a separate update event? Presumably
		// but whatever process consumes the QuestionnaireComplete message that this microservices pushes to SQS?
		err = dbConn.Create(&scheduledQuestionnaire)
		return

	case sql.ErrNoRows:
		// ad hoc
		err = ErrAdhocQuestionnaireCompleted
		return
	default:
		return
	}
}

func (event *QuestionnaireCompletedEvent) handleDeferFunc(err error, eventsQueue Queue, scheduledQuestionnaire *models.ScheduledQuestionnaire) {
	switch err {
	case nil:
		// pops the scheduled_questionnaire created message onto the events queue for asynchronous SQS transmission
		eventsQueue.Push(&ScheduledQuestionnaireEvent{
			Id:              scheduledQuestionnaire.Id,
			ParticipantId:   scheduledQuestionnaire.ParticipantId,
			QuestionnaireId: scheduledQuestionnaire.QuestionnaireId,
			Status:          scheduledQuestionnaire.Status.String,
			ScheduledAt:     scheduledQuestionnaire.ScheduledAt,
		})

	//	4. If not, push a new message to SQS that the user has completed all of their alloted scheduled questionnaires.
	// so, from this, I'm guessing the three scenarios for this would be if:
	//		- it's already completed
	//		- we've reached our maxiumum number of attempts
	// 		- it's adhoc and thus doesn't have/ require a scheduled questionnaire record
	case ErrMaxAttemptsReached, ErrScheduledQuestionnaireIsAlreadyCompleted, ErrAdhocQuestionnaireCompleted:
		eventsQueue.Push(event)

	default:
		// unexpected errors handled here, log and cry about it loudly!
		log.Fatalf("failed to process scheduled questionnaire event: %s", err)
	}
}
