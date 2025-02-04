-- +migrate Up
ALTER TABLE golang_crud.user
ADD COLUMN phone_number VARCHAR(20);