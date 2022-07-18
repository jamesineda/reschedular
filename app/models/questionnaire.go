package models

import "database/sql"

/*
	+----------------------+------------+----+---+-------+-----+
	|Field                 |Type        |Null|Key|Default|Extra|
	+----------------------+------------+----+---+-------+-----+
	|id                    |varchar(128)|NO  |PRI|NULL   |     |
	|study_id              |varchar(128)|NO  |   |NULL   |     |
	|name                  |varchar(128)|NO  |   |NULL   |     |
	|questions             |json        |NO  |   |NULL   |     |
	|max_attempts          |int(11)     |YES |   |NULL   |     |
	|hours_between_attempts|int(11)     |YES |   |24     |     |
	+----------------------+------------+----+---+-------+-----+
*/
type Questionnaire struct {
	Id                   string        `db:"id"`
	StudyId              string        `db:"study_id"`
	Name                 string        `db:"name"`
	Questions            string        `db:"questions"`
	MaxAttempts          sql.NullInt64 `db:"max_attempts"`
	HoursBetweenAttempts sql.NullInt64 `db:"max_attempts"`
}

type Questionnaires []*Questionnaire

// GetQuestions Parses JSON field to a struct (blocker: need to know the format of the JSON)
func (q *Questionnaire) GetQuestions() interface{} {
	panic("implement me")
}

// CanAttempt column will contain the maximum number of times a participant can fill in a given questionnaire
func (q *Questionnaire) CanAttempt(attemptNo int) bool {
	// if this is null, then there is no limit to the number of times they can fill it in.
	if !q.MaxAttempts.Valid {
		return true
	}

	return attemptNo < int(q.MaxAttempts.Int64)
}
