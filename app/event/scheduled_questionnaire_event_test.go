package event

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ScheduledQuestionnaireEventSuite struct {
	suite.Suite
}

func (suite *ScheduledQuestionnaireEventSuite) SetupTest() {
}

// TODO: finish this test
func (suite *ScheduledQuestionnaireEventSuite) Test_HandleEvent() {

}

func (suite *ScheduledQuestionnaireEventSuite) Test_ToSQSMessage() {

}

func TestScheduledQuestionnaireEventSuite(t *testing.T) {
	suite.Run(t, new(ScheduledQuestionnaireEventSuite))
}
