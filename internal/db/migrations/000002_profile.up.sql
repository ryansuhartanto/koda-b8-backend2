CREATE TABLE profiles (
    id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    name VARCHAR NOT NULL,
    picture_url VARCHAR
);

INSERT INTO profiles (id, name)
SELECT id, name FROM users;

ALTER TABLE users DROP COLUMN name;
