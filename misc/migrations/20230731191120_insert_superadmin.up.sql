INSERT INTO user_ (uuid, username, email, password, name, lastname, role, state)
SELECT
    '00000000-0000-0000-0000-000000000001',
    'pps_admin',
	'pps_admin@mail.com',
	'pps_admin',
	'admin',
	'admin',
	fcodes.role_code,
	fcodes.state_code
FROM (
    SELECT
        r.code as role_code,
        s.code as state_code
    FROM role r, user_state s
    WHERE r.description='super' AND s.description='activo'
) fcodes;
