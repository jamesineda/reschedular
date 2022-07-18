package event

import (
	"context"
	"github.com/jamesineda/reschedular/app/db"
	"github.com/jamesineda/reschedular/app/models"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ScheduledQuestionnaireEventSuite struct {
	suite.Suite
	PendingEvent   ScheduledQuestionnaireEvent
	CompletedEvent ScheduledQuestionnaireEvent
}

func (suite *ScheduledQuestionnaireEventSuite) SetupTest() {
	suite.PendingEvent = ScheduledQuestionnaireEvent{
		Name:            ScheduledQuestionnaire,
		Id:              "ABC123",
		ParticipantId:   "10",
		QuestionnaireId: "4",
		Status:          Pending,
	}

	suite.CompletedEvent = ScheduledQuestionnaireEvent{
		Name:            ScheduledQuestionnaire,
		Id:              "XYZ987",
		ParticipantId:   "10",
		QuestionnaireId: "5",
		Status:          Completed,
	}
}

// TODO: finish this test
func (suite *ScheduledQuestionnaireEventSuite) Test_HandleEvent() {
	suite.Run("when scheduled_questionnaire is complete", func() {
	})

	suite.Run("when questionnaire has reached max attempts", func() {
		fakeSQLX := db.NewSetFakeSQLX(nil, models.QuestionnaireResults{
			&models.QuestionnaireResult{},
			&models.QuestionnaireResult{},
			&models.QuestionnaireResult{},
		})

		conn, _ := db.NewFakeDatabaseConn(fakeSQLX)

		ctx := context.Background()
		context.WithValue(ctx, "db", conn)

	})

	suite.Run("when scheduled_questionnaire record can be created", func() {
	})
}

func TestScheduledQuestionnaireEventSuite(t *testing.T) {
	suite.Run(t, new(ScheduledQuestionnaireEventSuite))
}
