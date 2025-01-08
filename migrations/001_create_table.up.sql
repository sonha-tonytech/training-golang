-- +migrate Up
CREATE TABLE golang_crud.user (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_name VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE golang_crud.list (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    text VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    updated_at Timestamp NOT NULL
);
