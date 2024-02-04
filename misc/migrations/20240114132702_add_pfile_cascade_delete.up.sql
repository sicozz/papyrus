-- evidence
ALTER TABLE evidence
DROP CONSTRAINT IF EXISTS evidence_pfile_uuid_fkey;

ALTER TABLE evidence
ADD CONSTRAINT evidence_pfile_uuid_fkey
FOREIGN KEY (pfile_uuid)
REFERENCES pfile(uuid)
ON DELETE CASCADE;

-- approvation
ALTER TABLE approvation
DROP CONSTRAINT IF EXISTS approvation_pfile_uuid_fkey;

ALTER TABLE approvation
ADD CONSTRAINT approvation_pfile_uuid_fkey
FOREIGN KEY (pfile_uuid)
REFERENCES pfile(uuid)
ON DELETE CASCADE;

-- download
ALTER TABLE download
DROP CONSTRAINT IF EXISTS download_pfile_fkey;

ALTER TABLE download
ADD CONSTRAINT download_pfile_fkey
FOREIGN KEY (pfile)
REFERENCES pfile(uuid)
ON DELETE CASCADE;

-- upload
ALTER TABLE upload
DROP CONSTRAINT IF EXISTS upload_pfile_fkey;

ALTER TABLE upload
ADD CONSTRAINT upload_pfile_fkey
FOREIGN KEY (pfile)
REFERENCES pfile(uuid)
ON DELETE CASCADE;
