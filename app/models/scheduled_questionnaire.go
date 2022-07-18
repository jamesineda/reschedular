package models

import (
	"database/sql"
	"time"
)

/*
	+----------------+---------------------------+----+---+-------+-----+
	|Field           |Type                       |Null|Key|Default|Extra|
	+----------------+---------------------------+----+---+-------+-----+
	|id              |varchar(128)               |NO  |PRI|NULL   |     |
	|questionnaire_id|varchar(128)               |NO  |   |NULL   |     |
	|participant_id  |varchar(128)               |NO  |   |NULL   |     |
	|scheduled_at    |datetime                   |NO  |   |NULL   |     |
	|status          |enum('pending','completed')|YES |   |NULL   |     |
	+----------------+---------------------------+----+---+-------+-----+
*/
type ScheduledQuestionnaire struct {
	Id              string         `db:"id"`
	QuestionnaireId string         `db:"questionnaire_id"`
	ParticipantId   string         `db:"participant_id"`
	ScheduledAt     time.Time      `db:"scheduled_at"`
	Status          sql.NullString `db:"status"`
}

type ScheduledQuestionnaires []*ScheduledQuestionnaire
