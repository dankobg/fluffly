-- +goose Up
-- +goose StatementBegin
create type "gender" as enum ('m', 'f');

create table "user" (
  "id" uuid not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id")
);

create table "animal_type" (
  "id" bigint not null generated always as identity,
  "name" varchar(100) not null unique,
  primary key ("id")
);

-- create table "animal_breed" (
-- "id" bigint not null generated always as identity,
-- "animal_type_id" bigint not null references animal_types(id) on delete cascade,
-- "primary" varchar(255) not null,
-- "secondary" varchar(255) not null,
-- "mixed" boolean not null,
-- unique(animal_type_id, name),
-- primary key ("id")
-- );
create table "breed" (
  "id" bigint not null generated always as identity,
  "animal_type_id" bigint not null,
  "name" varchar(255) not null,
  primary key ("id"),
  constraint "uq_breed_animal_type_id_name_animal" unique ("animal_type_id", "name"),
  constraint "fk_breed_animal_type_id_animal" foreign key ("animal_type_id") references "animal_type" ("id") on delete cascade
);

create table "animal_species" (
  "id" bigint not null generated always as identity,
  "animal_type_id" bigint not null,
  "name" varchar(255) not null,
  primary key ("id"),
  constraint "uq_animal_species_animal_type_id_name_animal" unique ("animal_type_id", "name"),
  constraint "fk_animal_species_animal_type_id_animal" foreign key ("animal_type_id") references "animal_type" ("id") on delete cascade
);

create table "organization" (
  "id" bigint not null generated always as identity,
  "name" varchar(255) not null,
  "email" varchar(255) not null,
  "phone" varchar(50) not null,
  "address1" varchar(255),
  "address2" varchar(255),
  "city" varchar(255),
  "state" varchar(10),
  "postcode" varchar(20),
  "country" varchar(10),
  "url" text,
  "website" text,
  "mission_statement" text,
  "adoption_policy" text,
  "adoption_url" text,
  "distance" varchar(20),
  "facebook" text,
  "twitter" text,
  "youtube" text,
  "instagram" text,
  "pinterest" text,
  primary key ("id")
);

create table "organization_hour" (
  "id" bigint not null generated always as identity,
  "organization_id" bigint,
  "monday" text,
  "tuesday" text,
  "wednesday" text,
  "thursday" text,
  "friday" text,
  "saturday" text,
  "sunday" text,
  primary key ("id"),
  constraint fk_organization_hour_organization_id_organization foreign key ("organization_id") references "organization" ("id") on delete cascade
);

create table "organization_photo" (
  "id" bigint not null generated always as identity,
  "organization_id" bigint,
  "small" text,
  "medium" text,
  "large" text,
  "full" text,
  primary key ("id"),
  constraint "fk_organization_photo_organization_id_animal" foreign key ("organization_id") references "organization" ("id") on delete cascade
);

create table "animal" (
  "id" bigint not null generated always as identity,
  "user_id" uuid,
  "organization_id" bigint,
  "name" varchar not null,
  "type_id" bigint not null,
  "breed_id" bigint not null,
  "species_id" bigint not null,
  "gender" gender,
  "hermaphrodite" boolean not null default false,
  "age" varchar not null,
  "coat_length" varchar,
  "size" varchar not null,
  "image_url" varchar not null,
  "description" text,
  "adopted" boolean not null,
  "status" varchar(50),
  "status_changed_at" timestamptz,
  "distance" varchar,
  adopted_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  primary key ("id"),
  constraint "fk_animal_user_id_user" foreign key ("user_id") references "user" ("id") on delete cascade,
  constraint "fk_animal_organization_id_user" foreign key ("organization_id") references "organization" ("id") on delete set null

  -- -- attributes
  -- spayed_neutered boolean default false,
  -- house_trained boolean default false,
  -- declawed boolean default false,
  -- vaccinated boolean default false,
  -- special_needs boolean default false,
  -- "shots_current" boolean,
  -- good_with_kids boolean default true,
  -- good_with_cats boolean default true,
  -- good_with_dogs boolean default true,
  -- -- color
  -- -- contact
);

create table "animal_breed" (
  "animal_id" bigint,
  "breed_id" bigint,
  "primary" boolean not null default false,
  primary key ("animal_id", "breed_id"),
  constraint "fk_animal_breed_animal_id_animal" foreign key ("animal_id") references "animal" ("id") on delete cascade,
  constraint "fk_animal_breed_breed_id_animal" foreign key ("breed_id") references "breed" ("id") on delete cascade
);

create table user_animal_like (
  "user_id" uuid not null,
  "animal_id" bigint not null,
  "liked_at" timestamptz not null default now(),
  primary key ("user_id", "animal_id"),
  constraint "fk_user_animal_like_user_id_user" foreign key ("user_id") references "user" ("id") on delete cascade,
  constraint "fk_user_animal_like_animal_id_user" foreign key ("animal_id") references "animal" ("id") on delete cascade
);

create table "contact" (
  "id" bigint not null generated always as identity,
  "animal_id" bigint,
  "address_id" bigint not null,
  "phone" varchar not null,
  "email" varchar not null,
  "created_at" timestamptz not null default now(),
  "updated_at" timestamptz not null default now(),
  primary key ("id"),
  constraint "fk_contact_animal_id_animal" foreign key ("animal_id") references "animal" ("id") on delete cascade
);

create table "tag" (
  "id" bigint not null generated always as identity,
  "animal_id" bigint,
  "name" text,
  primary key ("id"),
  constraint "fk_tag_animal_id_animal" foreign key ("animal_id") references "animal" ("id") on delete cascade
);

create table "color" (
  "id" bigint not null generated always as identity,
  "animal_id" bigint,
  "color" varchar not null,
  primary key ("id"),
  constraint "fk_color_animal_id_animal" foreign key ("animal_id") references "animal" ("id") on delete cascade
);

create table "microchip" (
  "id" bigint not null generated always as identity,
  "animal_id" bigint,
  "number" varchar not null,
  "brand" varchar,
  "description" varchar,
  "location" varchar,
  primary key ("id"),
  constraint "uq_microchip_animal" unique ("animal_id"),
  constraint "fk_microchip_animal" foreign key ("animal_id") references "animal" ("id") on delete cascade
);

create table "animal_photo" (
  "id" bigint not null generated always as identity,
  "animal_id" bigint,
  "small" text,
  "medium" text,
  "large" text,
  "full" text,
  primary key ("id"),
  constraint "fk_animal_photo_animal_id_animal" foreign key ("animal_id") references "animal" ("id") on delete cascade
);

create table "animal_video" (
  "id" bigint not null generated always as identity,
  "animal_id" bigint,
  "url" text,
  primary key ("id"),
  constraint "fk_animal_video_animal_id_animal" foreign key ("animal_id") references "animal" ("id") on delete cascade
);

create table "adoption" (
  "id" bigint not null generated always as identity,
  "animal_id" bigint,
  "user_id" uuid,
  "adopted_at" timestamptz not null default now(),
  "returned_at" timestamptz,
  "notes" text,
  primary key ("id"),
  constraint "fk_adoption_animal_id_animal" foreign key ("animal_id") references "animal" ("id") on delete cascade,
  constraint "fk_adoption_user_id_animal" foreign key ("user_id") references "user" ("id") on delete set null
);

create table "address" (
  "id" bigint not null generated always as identity,
  "country_id" bigint,
  "unit_number" varchar(255),
  "street_number" varchar(255),
  "address_line_1" varchar(255),
  "address_line_2" varchar(255),
  "city" varchar(255),
  "region" varchar(255),
  "postal_code" varchar(20),
  "lat" decimal(9, 6),
  "lng" decimal(9, 6),
  "note" text,
  primary key ("id")
);

create table "country" (
  "id" bigint not null generated always as identity,
  "name" text not null,
  "iso_code" varchar(3) not null,
  primary key ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "user" cascade;
drop table if exists "animal_type" cascade;
drop table if exists "animal_breed" cascade;
drop table if exists "breed" cascade;
drop table if exists "animal_species" cascade;
drop table if exists "organization" cascade;
drop table if exists "organization_hour" cascade;
drop table if exists "organization_photo" cascade;
drop table if exists "animal" cascade;
drop table if exists "animal_breed" cascade;
drop table if exists "user_animal_like" cascade;
drop table if exists "contact" cascade;
drop table if exists "tag" cascade;
drop table if exists "color" cascade;
drop table if exists "microchip" cascade;
drop table if exists "animal_photo" cascade;
drop table if exists "animal_video" cascade;
drop table if exists "adoption" cascade;
drop table if exists "address" cascade;
drop table if exists "country" cascade;
drop type if exists "gender";
-- +goose StatementEnd
