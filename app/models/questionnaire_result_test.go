package models

import (
	"github.com/jamesineda/reschedular/app/utils"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type QuestionnaireResultTestSuite struct {
	suite.Suite
	Results QuestionnaireResults
}

func (suite *QuestionnaireResultTestSuite) SetupTest() {
	timer := utils.NewFakeTimer(time.Date(2022, 7, 18, 10, 0, 0, 0, time.UTC))

	now := timer.GetTimeNow()
	nowPlusOneHour := now.Add(1 * time.Hour)
	nowPlusTwoHours := now.Add(2 * time.Hour)
	nowPlusThreeHours := now.Add(3 * time.Hour)
	suite.Results = QuestionnaireResults{
		&QuestionnaireResult{Id: "ABC123", CompletedAt: &nowPlusTwoHours},
		&QuestionnaireResult{Id: "ABC456", CompletedAt: &now},
		&QuestionnaireResult{Id: "ABC789", CompletedAt: &nowPlusThreeHours},
		&QuestionnaireResult{Id: "XYZ987", CompletedAt: &nowPlusOneHour},
	}

}

func (suite *QuestionnaireResultTestSuite) Test_GetMostRecentResult() {
	suite.Run("get most recent result from collection that all have CompletedAt time", func() {
		suite.Equal(suite.Results[2], suite.Results.GetMostRecentResult())
	})

	suite.Run("get most recent result from collection that includes and incomplete result", func() {
		suite.Results[2].CompletedAt = nil
		suite.Equal(suite.Results[0], suite.Results.GetMostRecentResult())
	})
}

func TestQuestionnaireResult(t *testing.T) {
	suite.Run(t, new(QuestionnaireResultTestSuite))
}
