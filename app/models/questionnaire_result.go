package models

import (
	"database/sql"
	"time"
)

/*
	+-------------------------+------------+----+---+-------+-----+
	|Field                    |Type        |Null|Key|Default|Extra|
	+-------------------------+------------+----+---+-------+-----+
	|id                       |varchar(128)|NO  |   |NULL   |     |
	|answers                  |json        |NO  |   |NULL   |     |
	|questionnaire_id         |varchar(128)|NO  |   |NULL   |     |
	|participant_id           |varchar(128)|NO  |   |NULL   |     |
	|questionnaire_schedule_id|varchar(128)|YES |   |NULL   |     |
	|completed_at             |datetime    |YES |   |NULL   |     |
	+-------------------------+------------+----+---+-------+-----+
*/
type QuestionnaireResult struct {
	Id                      string         `db:"id"`
	Answers                 string         `db:"answers"`
	QuestionnaireId         string         `db:"questionnaire_id"`
	ParticipantId           string         `db:"participant_id"`
	QuestionnaireScheduleId sql.NullString `db:"questionnaire_schedule_id"`
	CompletedAt             *time.Time     `db:"completed_at"`
}

type QuestionnaireResults []*QuestionnaireResult

func (qr *QuestionnaireResults) Count() int {
	return len(*qr)
}

func (qr *QuestionnaireResults) GetMostRecentResult() (questionnaireResult *QuestionnaireResult) {
	for _, result := range *qr {
		if questionnaireResult == nil {
			questionnaireResult = result
		} else if result.CompletedAt != nil && result.CompletedAt.After(*questionnaireResult.CompletedAt) {
			questionnaireResult = result
		}
	}
	return
}
