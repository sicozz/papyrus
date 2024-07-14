ALTER TABLE plan
ADD date_check VARCHAR(2048);

ALTER TABLE task
ADD date_check VARCHAR(2048),
ADD date_close VARCHAR(2048);

ALTER TABLE pfile
ADD date_close VARCHAR(2048);
