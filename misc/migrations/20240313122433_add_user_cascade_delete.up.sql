ALTER TABLE permission
DROP CONSTRAINT IF EXISTS permission_user_uuid_fkey;

ALTER TABLE permission
ADD CONSTRAINT permission_user_uuid_fkey
FOREIGN KEY (user_uuid)
REFERENCES user_(uuid)
ON DELETE CASCADE;
