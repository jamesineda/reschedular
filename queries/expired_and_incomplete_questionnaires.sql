SELECT p.name AS participant_name, q.name AS questionnaire_name
FROM questionnaires q
         INNER JOIN scheduled_questionnaires sq ON q.id = sq.questionnaire_id
         INNER JOIN participants p ON sq.participant_id = p.id
WHERE sq.scheduled_at <= now()
  AND sq.status != 'completed';