CREATE TABLE evidence (
    task_uuid   UUID REFERENCES task(uuid),
    pfile_uuid  UUID REFERENCES pfile(uuid),
    PRIMARY KEY (task_uuid, pfile_uuid)
);
