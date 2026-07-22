ALTER TABLE users ADD COLUMN name VARCHAR;

UPDATE users
SET name = profiles.name
FROM profiles
WHERE profiles.id = users.id;

ALTER TABLE users ALTER COLUMN name SET NOT NULL;

DROP TABLE profiles;
