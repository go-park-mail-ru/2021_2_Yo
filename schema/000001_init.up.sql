CREATE TABLE users (
  id serial not null unique,
  name varchar(255) not null,
  surname varchar(255) not null,
  mail varchar(255) not null unique,
  password_hash varchar(255) not null
);

CREATE TABLE events (
  id serial not null unique,
  name varchar(255) not null,
  description varchar(255),
  views int
);