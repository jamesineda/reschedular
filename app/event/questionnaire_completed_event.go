package event

const (
	QuestionnaireCompleted = "QUESTIONNAIRE_COMPLETED"
)

type QuestionnaireCompletedEvent struct {
	Id                   string
	UserId               string
	StudyId              string
	QuestionnaireId      string
	CompletedAt          string
	RemainingCompletions int
}

func (q *QuestionnaireCompletedEvent) Name() string {
	return QuestionnaireCompleted
}


