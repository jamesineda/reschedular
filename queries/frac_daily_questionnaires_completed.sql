SELECT CAST(qr.completed_at AS date)                                          AS completed_on,
       SUM(CASE WHEN sq.status = 'completed' THEN 1 ELSE 0 END) / COUNT(q.id) AS fraction_complete
FROM questionnaires q
         INNER JOIN scheduled_questionnaires sq ON q.id = sq.questionnaire_id
         INNER JOIN questionnaire_results qr ON q.id = qr.questionnaire_id
GROUP BY 1;