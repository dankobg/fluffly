-- +goose Up
-- +goose StatementBegin
create type "gender" as enum ('m', 'f');

create table "user" (
  "id" uuid not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id")
);

create table "country" (
  "id" bigint not null generated always as identity,
  "name" text not null,
  "iso_alpha2" varchar(2) not null,
  "iso_alpha3" varchar(3) not null,
  "iso_numeric" varchar(3) not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "uq_country_name" unique ("name"),
  constraint "uq_country_iso_alpha2" unique ("iso_alpha2"),
  constraint "uq_country_iso_alpha3" unique ("iso_alpha3")
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
  "postal_code" varchar(50),
  "lat" decimal(9, 6),
  "lng" decimal(9, 6),
  "note" text,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "fk_address_country_id_country" foreign key ("country_id") references "country" ("id") on delete cascade
);

create table "contact" (
  "id" bigint not null generated always as identity,
  "user_id" uuid,
  "organization_id" bigint,
  "address_id" bigint not null,
  "phone" varchar(30) not null,
  "email" text not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "fk_contact_address_id_animal" foreign key ("address_id") references "address" ("id") on delete cascade,
  constraint "fk_contact_user_or_organization_present" check (
    (user_id is not null and organization_id is null) or
    (user_id is null and organization_id is not null)
  )
);

create table "animal_type" (
  "id" bigint not null generated always as identity,
  "name" varchar(100) not null unique,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id")
);

create table "breed" (
  "id" bigint not null generated always as identity,
  "animal_type_id" bigint not null,
  "name" varchar(255) not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "uq_breed_animal_type_id_name_animal" unique ("animal_type_id", "name"),
  constraint "fk_breed_animal_type_id_animal" foreign key ("animal_type_id") references "animal_type" ("id") on delete cascade
);

create table "animal_species" (
  "id" bigint not null generated always as identity,
  "animal_type_id" bigint not null,
  "name" varchar(255) not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "uq_animal_species_animal_type_id_name_animal" unique ("animal_type_id", "name"),
  constraint "fk_animal_species_animal_type_id_animal" foreign key ("animal_type_id") references "animal_type" ("id") on delete cascade
);

create table "organization" (
  "id" bigint not null generated always as identity,
  "contact_id" bigint not null,
  "name" varchar(255) not null,
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
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "fk_organization_contact_id_contact" foreign key ("contact_id") references "contact" ("id") on delete cascade
);

create table "organization_hour" (
  "id" bigint not null generated always as identity,
  "organization_id" bigint,
  "monday" varchar(30),
  "tuesday" varchar(30),
  "wednesday" varchar(30),
  "thursday" varchar(30),
  "friday" varchar(30),
  "saturday" varchar(30),
  "sunday" varchar(30),
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
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
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "fk_organization_photo_organization_id_animal" foreign key ("organization_id") references "organization" ("id") on delete cascade
);

create table "animal" (
  "id" bigint not null generated always as identity,
  "contact_id" bigint not null,
  "type_id" bigint not null,
  "breed_id" bigint not null,
  "species_id" bigint not null,
  "name" varchar(255) not null,
  "gender" gender,
  "hermaphrodite" boolean not null default false,
  "age" varchar(20) not null,
  "coat_length" varchar(30),
  "size" varchar(30) not null,
  "image_url" text not null,
  "description" text,
  "distance" varchar,
  "status" varchar(30),
  "status_changed_at" timestamptz,
  "adopted_at" timestamptz,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "fk_animal_contact_id_contact" foreign key ("contact_id") references "contact" ("id") on delete cascade
  
  -- @TODO: handle properties unique for each species
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
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
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

create table "tag" (
  "id" bigint not null generated always as identity,
  "animal_id" bigint,
  "name" varchar(255),
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "fk_tag_animal_id_animal" foreign key ("animal_id") references "animal" ("id") on delete cascade
);

create table "color" (
  "id" bigint not null generated always as identity,
  "animal_id" bigint,
  "color" varchar(255) not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "fk_color_animal_id_animal" foreign key ("animal_id") references "animal" ("id") on delete cascade
);

create table "microchip" (
  "id" bigint not null generated always as identity,
  "animal_id" bigint,
  "number" varchar(30) not null,
  "brand" varchar(255),
  "description" text,
  "location" varchar(100),
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
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
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "fk_animal_photo_animal_id_animal" foreign key ("animal_id") references "animal" ("id") on delete cascade
);

create table "animal_video" (
  "id" bigint not null generated always as identity,
  "animal_id" bigint,
  "url" text,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
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
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "fk_adoption_animal_id_animal" foreign key ("animal_id") references "animal" ("id") on delete cascade,
  constraint "fk_adoption_user_id_animal" foreign key ("user_id") references "user" ("id") on delete set null
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
