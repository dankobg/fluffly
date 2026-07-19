-- +goose Up
-- +goose StatementBegin
-- insert cat species
insert into "animal_specie" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values
    ('Cat')
) as sp("name")
where at."name" = 'Cat';


-- insert dog species
insert into "animal_specie" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values
  ('Dog')
) as sp("name")
where at."name" = 'Dog';


-- insert rabbit species
insert into "animal_specie" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values
  ('Rabbit')
) as sp("name")
where at."name" = 'Rabbit';


-- insert horse species
insert into "animal_specie" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values
    ('Donkey'),
    ('Horse'),
    ('Mule'),
    ('Pony'),
    ('Miniature Horse')
) as sp("name")
where at."name" = 'Horse';


-- insert bird species
insert into "animal_specie" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values
    ('Button-Quail'),
    ('Chicken'),
    ('Dove'),
    ('Duck'),
    ('Emu'),
    ('Finch'),
    ('Goose'),
    ('Guinea Fowl'),
    ('Ostritch'),
    ('Parakeet'),
    ('Parrot'),
    ('Peacock / Peafowl'),
    ('Pheasant'),
    ('Quail'),
    ('Rhea'),
    ('Swan'),
    ('Toucan'),
    ('Turkey')
) as sp("name")
where at."name" = 'Bird';


-- insert barnyard species
insert into "animal_specie" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values
    ('Alpaca'),
    ('Cow'),
    ('Goat'),
    ('Llama'),
    ('Pig'),
    ('Pot Bellied'),
    ('Sheep')
) as sp("name")
where at."name" = 'Barnyard';


-- insert small and furry species
insert into "animal_specie" ("animal_type_id", "name")
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
    ('Degu'),
    ('Prairie Dog'),
    ('Skunk'),
    ('Sugar Glider')
) as sp("name")
where at."name" = 'Small & Furry';


-- insert scales fins and other species
insert into "animal_specie" ("animal_type_id", "name")
select at."id", sp."name"
from "animal_type" at
cross join (
  values
    ('Amphibian'),
    ('Fish'),
    ('Other Animal'),
    ('Reptile'),
    ('Salamander / Newt'),
    ('Snake'),
    ('Toad'),
    ('Tortoise'),
    ('Turtle')
) as sp("name")
where at."name" = 'Scales, Fins & Other';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate table "animal_specie" cascade;
-- +goose StatementEnd
