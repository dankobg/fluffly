-- +goose Up
-- +goose StatementBegin
insert into "animal_type" (name) values 
  ('Dog'),
  ('Cat'),
  ('Rabbit'),
  ('Horse'),
  ('Bird'),
  ('Barnyard'),
  ('Small & Furry'),
  ('Scales, Fins & Other');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate table "animal_type" cascade;
-- +goose StatementEnd
