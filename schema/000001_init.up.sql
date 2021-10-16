CREATE TABLE "user" (
  id serial not null unique,
  name varchar(255) not null,
  surname varchar(255) not null,
  mail varchar(255) not null unique,
  password varchar(255) not null
);

CREATE TABLE "event" (
  id serial not null unique,
  name varchar(255) not null,
  description varchar(255),
  views int
);