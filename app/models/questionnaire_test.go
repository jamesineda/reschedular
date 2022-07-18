package models

import (
	"database/sql"
	"github.com/stretchr/testify/suite"
	"testing"
)

type QuestionnaireTestSuite struct {
	suite.Suite
}

func (suite *QuestionnaireTestSuite) SetupTest() {

}

func (suite *QuestionnaireTestSuite) Test_CanAttempt() {
	suite.Run("when max attempts is null", func() {
		questionnaire := &Questionnaire{}
		suite.Equal(true, questionnaire.CanAttempt(1000))
	})

	suite.Run("when max attempts is 3", func() {
		questionnaire := &Questionnaire{MaxAttempts: sql.NullInt64{Int64: 3, Valid: true}}
		suite.Run("and attempt no 4 is specified", func() {
			suite.Equal(false, questionnaire.CanAttempt(4))
		})

		suite.Run("and attempt no 2 is specified", func() {
			suite.Equal(true, questionnaire.CanAttempt(2))
		})
	})
}

func (suite *QuestionnaireTestSuite) Test_GetHoursBetweenAttemptsDuration() {
	suite.Run("attempts between hours is specified as 4", func() {
		questionnaire := Questionnaire{HoursBetweenAttempts: sql.NullInt64{Int64: 4, Valid: true}}
		suite.Equal(float64(4), questionnaire.GetHoursBetweenAttemptsDuration().Hours())

	})

	suite.Run("attempts between hours is not specified", func() {
		questionnaire := Questionnaire{}
		suite.Equal(float64(24), questionnaire.GetHoursBetweenAttemptsDuration().Hours())
	})
}

func TestQuestionnaire(t *testing.T) {
	suite.Run(t, new(QuestionnaireTestSuite))
}
