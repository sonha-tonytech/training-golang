-- +migrate Down
ALTER TABLE golang_crud.user
DROP COLUMN phone_number;