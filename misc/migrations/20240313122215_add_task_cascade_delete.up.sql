ALTER TABLE evidence
DROP CONSTRAINT IF EXISTS evidence_task_uuid_fkey;

ALTER TABLE evidence
ADD CONSTRAINT evidence_task_uuid_fkey
FOREIGN KEY (task_uuid)
REFERENCES task(uuid)
ON DELETE CASCADE;
