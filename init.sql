CREATE DATABASE papyrus;

\connect papyrus;

CREATE TABLE dir (
    uuid        UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(256)  NOT NULL,
    parent_dir  UUID          REFERENCES dir NOT NULL,
    CONSTRAINT different_uuid_parent_dir CHECK ((uuid <> parent_dir) OR uuid = '00000000-0000-0000-0000-000000000000')
);

CREATE TABLE role (
    code         SERIAL       PRIMARY KEY,
    description  VARCHAR(32)  UNIQUE NOT NULL
);

CREATE TABLE user_state (
    code         SERIAL       PRIMARY KEY,
    description  VARCHAR(32)  UNIQUE NOT NULL
);

CREATE TABLE user_ (
    uuid      UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    username  VARCHAR(32)  NOT NULL,
    email     VARCHAR(64)  UNIQUE NOT NULL,
    password  VARCHAR(32)  NOT NULL,
    name      VARCHAR(32)  NOT NULL,
    lastname  VARCHAR(32)  NOT NULL,
    role      SERIAL       REFERENCES role NOT NULL,
    state     SERIAL       REFERENCES user_state NOT NULL
);

CREATE TABLE pfile_type (
    code         SERIAL       PRIMARY KEY,
    description  VARCHAR(32)  UNIQUE NOT NULL
);

CREATE TABLE pfile_state (
    code         SERIAL       PRIMARY KEY,
    description  VARCHAR(32)  UNIQUE NOT NULL
);

CREATE TABLE pfile (
    uuid           UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    code           VARCHAR(1024)   NOT NULL,
    name           VARCHAR(1024)  NOT NULL,
    fs_path        VARCHAR(1024)  NOT NULL,
    date_creation  TIMESTAMP     NOT NULL,
    date_input     TIMESTAMP     NOT NULL,
    type           SERIAL        REFERENCES pfile_type NOT NULL,
    state          SERIAL        REFERENCES pfile_state NOT NULL,
    dir            UUID          REFERENCES dir NOT NULL,
    version        VARCHAR(32)   NOT NULL,
    term           INTEGER       NOT NULL,
    subtype        VARCHAR(32)   NOT NULL,
    resp_user      UUID          REFERENCES user_(uuid)
);

CREATE TABLE version (
    uuid    UUID       primary key default gen_random_uuid(),
    date    TIMESTAMP  NOT NULL,
    pfile   UUID       REFERENCES pfile NOT NULL
);

CREATE TABLE approvation (
    user_uuid   UUID  REFERENCES user_(uuid),
    pfile_uuid  UUID  REFERENCES pfile(uuid),
    is_approved BOOLEAN NOT NULL,
    PRIMARY KEY (user_uuid, pfile_uuid)
);

CREATE TABLE download (
    uuid   UUID       PRIMARY KEY DEFAULT gen_random_uuid(),
    date   TIMESTAMP  NOT NULL,
    user_  UUID       REFERENCES user_ NOT NULL,
    pfile  UUID       REFERENCES pfile NOT NULL
);

CREATE TABLE upload (
    uuid   UUID       PRIMARY KEY DEFAULT gen_random_uuid(),
    date   TIMESTAMP  NOT NULL,
    user_  UUID       REFERENCES user_ NOT NULL,
    pfile   UUID       REFERENCES pfile NOT NULL
);

CREATE TABLE permission (
    user_uuid   UUID  REFERENCES user_(uuid),
    dir_uuid    UUID  REFERENCES dir(uuid),
    PRIMARY KEY (user_uuid, dir_uuid)
);

CREATE TABLE project_state (
    code         SERIAL       PRIMARY KEY,
    description  VARCHAR(32)  UNIQUE NOT NULL
);

CREATE TABLE project (
    uuid         UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(1024)    NOT NULL,
    description  VARCHAR(2048)  NOT NULL,
    state        SERIAL         REFERENCES PROJECT_STATE NOT NULL,
    dir          UUID           REFERENCES DIR NOT NULL
);

CREATE TABLE plan_state (
    code         SERIAL       PRIMARY KEY,
    description  VARCHAR(32)  UNIQUE NOT NULL
);

CREATE TABLE plan (
    uuid              UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    title             VARCHAR(1024)  NOT NULL,
    description       VARCHAR(2048)  NOT NULL,
    origin            VARCHAR(2048)  NOT NULL,
    analysis          VARCHAR(2048)  NOT NULL,
    discovery_date    TIMESTAMP      NOT NULL,
    record_date       TIMESTAMP      NOT NULL,
    termination_date  TIMESTAMP      NOT NULL,
    state             SERIAL         REFERENCES plan_state NOT NULL,
    stage             SERIAL         NOT NULL,
    project           UUID           REFERENCES project NOT NULL,
    issuing_user_     UUID           REFERENCES user_ NOT NULL,
    offender_user_    UUID           REFERENCES user_ NOT NULL,
    assigned_user_    UUID           REFERENCES user_
);

CREATE TABLE task_state (
    code         SERIAL       PRIMARY KEY,
    description  VARCHAR(32)  UNIQUE NOT NULL
);

CREATE TABLE task (
    uuid           UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(1024)    NOT NULL,
    procedure      VARCHAR(2048)  NOT NULL,
    date_creation  TIMESTAMP      NOT NULL,
    term           INTEGER        NOT NULL,
    -- deadline       TIMESTAMP      NOT NULL,
    state          SERIAL         REFERENCES task_state NOT NULL,
    dir            UUID           REFERENCES dir NOT NULL,
    -- evidence_dir   UUID           REFERENCES dir NOT NULL,
    creator_user   UUID           REFERENCES user_ NOT NULL,
    recv_user      UUID           REFERENCES user_,
    chk            BOOLEAN        NOT NULL
    -- plan           UUID           REFERENCES plan
);
