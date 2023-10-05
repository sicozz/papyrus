INSERT INTO role (description)
VALUES ('estandar'), ('admin'), ('super');

INSERT INTO user_state (description)
VALUES ('inactivo'), ('activo');

INSERT INTO pfile_type (description)
VALUES ('documento'), ('formato'), ('registro');

INSERT INTO pfile_state (description)
VALUES ('revision'), ('activo'), ('obsoleto');

INSERT INTO project_state (description)
VALUES ('inactivo'), ('activo'), ('cerrado');

INSERT INTO plan_state (description)
VALUES ('abierto'), ('cerrado'), ('abandonado');

INSERT INTO task_state (description)
VALUES ('abierta'), ('cerrada'), ('abandonada');
