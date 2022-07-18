package models

/*
	+-----+------------+----+---+-------+-----+
	|Field|Type        |Null|Key|Default|Extra|
	+-----+------------+----+---+-------+-----+
	|id   |varchar(128)|NO  |PRI|NULL   |     |
	|name |varchar(128)|NO  |   |NULL   |     |
	+-----+------------+----+---+-------+-----+

*/
type Participant struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}

type Participants []*Participant
