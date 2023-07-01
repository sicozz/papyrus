INSERT INTO role (description)
VALUES ('estandar'), ('admin'), ('super');

INSERT INTO user_state (description)
VALUES ('inactivo'), ('activo');

INSERT INTO file_type (description)
VALUES ('documento'), ('formato');

INSERT INTO file_state (description)
VALUES ('inactivo'), ('activo'), ('obsoleto');

INSERT INTO file_stage (description)
VALUES ('cargado'), ('revisado'), ('aprobado');

INSERT INTO project_state (description)
VALUES ('inactivo'), ('activo'), ('cerrado');

INSERT INTO plan_state (description)
VALUES ('abierto'), ('cerrado'), ('abandonado');

INSERT INTO task_state (description)
VALUES ('abierta'), ('cerrada'), ('abandonada');
