-- +goose Up
CREATE TABLE persons (
	id	INTEGER NOT NULL UNIQUE,
	username text UNIQUE,
	age	text,
	name text,
	family	text,
	password	text,
	role	text,
	PRIMARY KEY(id AUTOINCREMENT)
);

-- +goose Down
DROP TABLE persons;
