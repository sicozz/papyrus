-- Drop columns from the 'plan' table
ALTER TABLE plan
DROP COLUMN IF EXISTS title,
DROP COLUMN IF EXISTS description,
DROP COLUMN IF EXISTS origin,
DROP COLUMN IF EXISTS state,
DROP COLUMN IF EXISTS stage,
DROP COLUMN IF EXISTS analysis,
DROP COLUMN IF EXISTS discovery_date,
DROP COLUMN IF EXISTS record_date,
DROP COLUMN IF EXISTS termination_date,
DROP COLUMN IF EXISTS project,
DROP COLUMN IF EXISTS issuing_user_,  -- You have an extra underscore here
DROP COLUMN IF EXISTS offender_user_,
DROP COLUMN IF EXISTS assigned_user_;

-- Add columns to the 'plan' table
ALTER TABLE plan
ADD code              VARCHAR(1024)  UNIQUE NOT NULL,
ADD name              VARCHAR(2048)  UNIQUE NOT NULL,
ADD origin            VARCHAR(2048)  NOT NULL,
ADD action_type       VARCHAR(2048)  NOT NULL,
ADD term              INTEGER        NOT NULL,
ADD creator_user      UUID           REFERENCES user_ NOT NULL,
ADD resp_user         UUID           REFERENCES user_ NOT NULL,
ADD date_create       TIMESTAMP      NOT NULL,
ADD date_close        TIMESTAMP      NOT NULL,
ADD causes            VARCHAR(2048)  NOT NULL,
ADD conclusions       VARCHAR(2048)  NOT NULL,
ADD state             VARCHAR(2048)  NOT NULL,
ADD stage             SERIAL         NOT NULL,
ADD dir               UUID           REFERENCES dir NOT NULL,
ADD action0_desc      VARCHAR(2048),
ADD action0_date      VARCHAR(2048),
ADD action0_user      VARCHAR(36),
ADD action1_desc      VARCHAR(2048),
ADD action1_date      VARCHAR(2048),
ADD action1_user      VARCHAR(36),
ADD action2_desc      VARCHAR(2048),
ADD action2_date      VARCHAR(2048),
ADD action2_user      VARCHAR(36),
ADD action3_desc      VARCHAR(2048),
ADD action3_date      VARCHAR(2048),
ADD action3_user      VARCHAR(36),
ADD action4_desc      VARCHAR(2048),
ADD action4_date      VARCHAR(2048),
ADD action4_user      VARCHAR(36),
ADD action5_desc      VARCHAR(2048),
ADD action5_date      VARCHAR(2048),
ADD action5_user      VARCHAR(36);

-- Create the 'action' table
-- CREATE TABLE action (
--     uuid          UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
--     description   VARCHAR(2048)  NOT NULL,
--     creator_user  UUID           REFERENCES user_ NOT NULL,
--     date_close    TIMESTAMP      NOT NULL  -- Removed semicolon here
-- );
