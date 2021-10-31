CREATE TABLE "user" (
                        id serial not null unique,
                        name varchar(50) not null,
                        surname varchar(50) not null,
                        mail varchar(150) not null unique,
                        password varchar(50) not null,
                        about varchar(150) not null
);

CREATE TABLE "event" (
                         id serial not null unique,
                         title varchar(255) not null,
                         description varchar(500) not null,
                         text varchar(1000) not null,
                         city varchar(255) not null,
                         category varchar(255) not null,
                         viewed BIGINT not null,
                         img_url varchar(500) not null,
                         date varchar(10) not null,
                         geo varchar(255) not null,
                         author_id int references "user" (id) on delete cascade not null
);

CREATE TABLE "tag" (
                       id serial not null unique,
                       Name varchar(30) not null
);

CREATE TABLE "tag_event" (
                             ID serial not null unique,
                             tag_id int references "tag" (id) on delete cascade not null,
                             event_id int references "event" (id) on delete cascade not null
);