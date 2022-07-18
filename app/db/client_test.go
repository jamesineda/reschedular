package db

import (
	"github.com/jamesineda/reschedular/app/models"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ClientTestSuite struct {
	suite.Suite
}

func (suite *ClientTestSuite) SetupTest() {

}

func (suite *ClientTestSuite) Test_generateSelectQuery() {
	suite.Run("generate query for GetList", func() {

		table := models.QuestionnaireResults{}
		tableName, selectFields := getSelectOptions(table)
		filters := Filters{
			[]interface{}{"name", "=", "hair regrowth questionnaire"},
			[]interface{}{"max_attempts", "=", 5},
			[]interface{}{"questions", "=", `{"did your hair grow back?": "no"}`},
		}

		query := generateSelectQuery(tableName, selectFields, filters)
		suite.Equal(`SELECT id,answers,questionnaire_id,participant_id,questionnaire_schedule_id,completed_at FROM questionnaire_results WHERE name = ? AND max_attempts = ? AND questions = ?`, query)
	})
}

func (suite *ClientTestSuite) Test_Filter_Values() {
	suite.Run("generate Filter values", func() {
		filters := Filters{
			[]interface{}{"name", "=", "hair regrowth questionnaire"},
			[]interface{}{"max_attempts", "=", 5},
			[]interface{}{"questions", "=", `{"did your hair grow back?": "no"}`},
		}

		values := filters.Values()
		suite.Equal([]interface{}{"hair regrowth questionnaire", 5, `{"did your hair grow back?": "no"}`}, values)
	})
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
