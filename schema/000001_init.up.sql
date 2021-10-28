CREATE TABLE "user" (
  ID serial not null unique,
  user_id int not null,
  Name varchar(255) not null,
  Surname varchar(255) not null,
  Mail varchar(255) not null unique,
  Password varchar(255) not null
);

CREATE TABLE "event" (
  ID serial not null unique,
  event_id int not null,
  Title varchar(255) not null,
  Description varchar(255),
  Text varchar(1000),
  City varchar(255) not null,
  Category varchar(255),
  Viewed BIGINT not null,
  ImgUrl varchar(500) not null,
  Date varchar(255) not null,
  GEO varchar(255) not null,
  Author_id int references "user" (ID) on delete cascade not null
);

CREATE TABLE "tag" (
  ID serial not null unique,
  Name varchar(255) not null
);

CREATE TABLE "tag_event" (
  ID serial not null unique,
  tag_id int references "tag" (ID) on delete cascade not null,
  event_id int references "event" (ID) on delete cascade not null
);