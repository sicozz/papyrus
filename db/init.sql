CREATE DATABASE papyrus;

\connect papyrus;

CREATE TABLE DIR (
    uuid        UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(256)  NOT NULL,
    parent_dir  UUID          REFERENCES dir
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
    username  VARCHAR(32)  UNIQUE NOT NULL,
    email     VARCHAR(64)  UNIQUE NOT NULL,
    password  VARCHAR(32)  NOT NULL,
    name      VARCHAR(32)  NOT NULL,
    lastname  VARCHAR(32)  NOT NULL,
    role      SERIAL       REFERENCES role NOT NULL,
    state     SERIAL       REFERENCES user_state NOT NULL
);

CREATE TABLE file_type (
    code         SERIAL       PRIMARY KEY,
    description  VARCHAR(32)  UNIQUE NOT NULL
);

CREATE TABLE file_state (
    code         SERIAL       PRIMARY KEY,
    description  VARCHAR(32)  UNIQUE NOT NULL
);

CREATE TABLE file_stage (
    code         SERIAL       PRIMARY KEY,
    description  VARCHAR(32)  UNIQUE NOT NULL
);

CREATE TABLE file (
    uuid           UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    code           VARCHAR(32)   NOT NULL,
    path           VARCHAR(256)  NOT NULL,
    creation_date  TIMESTAMP     NOT NULL,
    input_date     TIMESTAMP     NOT NULL,
    type           SERIAL        REFERENCES FILE_TYPE NOT NULL,
    state          SERIAL        REFERENCES FILE_STATE NOT NULL,
    stage          SERIAL        REFERENCES FILE_STAGE NOT NULL,
    dir            UUID          REFERENCES dir NOT NULL,
    revision_user  UUID          REFERENCES user_ NOT NULL,
    approval_user  UUID          REFERENCES user_ NOT NULL
);

create table version (
    uuid  uuid       primary key default gen_random_uuid(),
    date  TIMESTAMP  NOT NULL,
    file  UUID       REFERENCES file NOT NULL
);

CREATE TABLE download (
    uuid   UUID       PRIMARY KEY DEFAULT gen_random_uuid(),
    date   TIMESTAMP  NOT NULL,
    user_  UUID       REFERENCES user_ NOT NULL,
    file   UUID       REFERENCES file NOT NULL
);

CREATE TABLE upload (
    uuid   UUID       PRIMARY KEY DEFAULT gen_random_uuid(),
    date   TIMESTAMP  NOT NULL,
    user_  UUID       REFERENCES user_ NOT NULL,
    file   UUID       REFERENCES file NOT NULL
);

CREATE TABLE read_permission (
    uuid     UUID     PRIMARY KEY DEFAULT gen_random_uuid(),
    allowed  BOOLEAN  NOT NULL,
    user_    UUID     REFERENCES user_ NOT NULL,
    file     UUID     REFERENCES file NOT NULL
);

CREATE TABLE write_permission (
    uuid     UUID     PRIMARY KEY DEFAULT gen_random_uuid(),
    allowed  BOOLEAN  NOT NULL,
    user_    UUID     REFERENCES user_ NOT NULL,
    file     UUID     REFERENCES file NOT NULL
);

CREATE TABLE project_state (
    code         SERIAL       PRIMARY KEY,
    description  VARCHAR(32)  UNIQUE NOT NULL
);

CREATE TABLE project (
    uuid         UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(64)    NOT NULL,
    description  VARCHAR(1024)  NOT NULL,
    state        SERIAL         REFERENCES PROJECT_STATE NOT NULL,
    dir          UUID           REFERENCES DIR NOT NULL
);

CREATE TABLE plan_state (
    code         SERIAL       PRIMARY KEY,
    description  VARCHAR(32)  UNIQUE NOT NULL
);

CREATE TABLE plan (
    uuid              UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    title             VARCHAR(64)    NOT NULL,
    description       VARCHAR(1024)  NOT NULL,
    origin            VARCHAR(1024)  NOT NULL,
    analysis          VARCHAR(1024)  NOT NULL,
    discovery_date    TIMESTAMP      NOT NULL,
    record_date       TIMESTAMP      NOT NULL,
    termination_date  TIMESTAMP      NOT NULL,
    state             SERIAL         REFERENCES plan_state NOT NULL,
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
    title          VARCHAR(64)    NOT NULL,
    description    VARCHAR(1024)  NOT NULL,
    date           TIMESTAMP      NOT NULL,
    deadline       TIMESTAMP      NOT NULL,
    state          SERIAL         REFERENCES task_state NOT NULL,
    dir            UUID           REFERENCES dir NOT NULL,
    evidence_dir   UUID           REFERENCES dir NOT NULL,
    issuing_user   UUID           REFERENCES user_ NOT NULL,
    assigned_user  UUID           REFERENCES user_,
    plan           UUID           REFERENCES plan
);
