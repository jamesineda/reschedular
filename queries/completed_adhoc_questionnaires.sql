SELECT p.name AS participant_name, q.name AS questionnaire_name, qr.completed_at AS ad_hoc_timestamp
FROM participants p
    LEFT JOIN scheduled_questionnaires sq
ON p.id = sq.participant_id
    INNER JOIN questionnaire_results qr ON p.id = qr.participant_id
    INNER JOIN questionnaires q ON qr.questionnaire_id = q.id
-- checking that qr.completed_at is not NULL to account for partially completely questionnaires; not sure if that's a thing?
-- I'm assuming questionnaires done outside of the scheduled designed are the 'ad hoc' questionnaires
WHERE sq.participant_id IS NULL AND qr.questionnaire_schedule_id IS NULL AND qr.completed_at IS NOT NULL;