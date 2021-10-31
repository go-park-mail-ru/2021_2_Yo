CREATE TABLE "user" (
                        id serial not null unique,
                        name varchar(50) not null,
                        surname varchar(50) not null,
                        mail varchar(150) not null unique,
                        password varchar(255) not null,
                        about varchar(150)
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

INSERT INTO "user" (name, surname, mail, password, about)
VALUES ('Andrey', 'Ivanov', 'test@mail.ru', 'hashhashhash', 'I am soooo coool DAAAMN'),
       ('Ivan', 'Andreev', 'kool@mail.ru', 'hashhashhash', 'My name is Ivan Andreev'),
       ('Petr', 'Leshiy', 'clown@mail.ru', 'hashhashhash', 'Yo my desctiption');

INSERT INTO "event" (title, description, text, city, category, viewed, img_url, date, geo, author_id)
VALUES ('Tusovka', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina', 1),
       ('Funny party', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina', 3),
       ('Really funny party', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina', 2),
       ('PAPAPAPA', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina', 1),
       ('Joked party', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina', 2),
       ('Boring party', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina', 3),
       ('funny sunny', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina', 2),
       ('Test party', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina', 1);

INSERT INTO "tag" (name)
VALUES ('Party'),
       ('Alcohol'),
       ('Funny'),
       ('Boring'),
       ('Popular');

INSERT INTO "tag_event" (tag_id, event_id)
VALUES (1, 1),
       (1, 2),
       (2, 1),
       (2, 4),
       (5, 1),
       (5, 2),
       (3, 1),
       (3, 6),
       (3, 5);