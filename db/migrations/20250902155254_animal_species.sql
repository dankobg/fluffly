-- +goose Up
-- +goose StatementBegin
-- insert barnyard breeds
-- insert cat species
insert into "animal_species" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values ('Cat')
) as sp("name")
where at."name" = 'Cat';


-- insert dog species
insert into "animal_species" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values
    ('Dog')
) as sp("name")
where at."name" = 'Dog';


-- insert rabbit species
insert into "animal_species" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values ('Rabbit')
) as sp("name")
where at."name" = 'Rabbit';


-- insert horse species
insert into "animal_species" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values ('Horse')
) as sp("name")
where at."name" = 'Horse';


-- insert bird species
insert into "animal_species" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values ('Bird')
) as sp("name")
where at."name" = 'Bird';


-- insert barnyard species
insert into "animal_species" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values ('Barnyard')
) as sp("name")
where at."name" = 'Barnyard';


-- insert small and furry species
insert into "animal_species" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values
  ('Guinea Pig'),
  ('Chinchilla'),
  ('Rat'),
  ('Mouse'),
  ('Ferret'),
  ('Hamster'),
  ('Hedgehog'),
  ('Gerbil'),
  ('Degu')
) as sp("name")
where at."name" = 'Small & Furry';


-- insert scales fins and other species
insert into "animal_species" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values
    ('Reptile'),
    ('Other Animal'),
    ('Snake'),
    ('Tortoise'),
    ('Turtl')
) as sp("name")
where at."name" = 'Scales, Fins & Other';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate table "animal_species" cascade;
-- +goose StatementEnd
