package event

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type QuestionnaireCompletedEventSuite struct {
	suite.Suite
}

func (suite *QuestionnaireCompletedEventSuite) SetupTest() {
}

// TODO: finish this test
func (suite *QuestionnaireCompletedEventSuite) Test_HandleEvent() {
	suite.Run("when there are no remaining completions", func() {
	})

	suite.Run("when questionnaire has reached max attempts", func() {

	})

	suite.Run("when there are no associated scheduled_questionnaire records", func() {
	})

	suite.Run("when there IS remaining completions", func() {
	})
}

func (suite *QuestionnaireCompletedEventSuite) Test_ToSQSMessage() {

}

func TestQuestionnaireCompletedEventSuite(t *testing.T) {
	suite.Run(t, new(QuestionnaireCompletedEventSuite))
}
