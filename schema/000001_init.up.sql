CREATE TABLE "user" (
                        id serial not null unique,
                        name varchar(50) not null,
                        surname varchar(50) not null,
                        about varchar(150),
                        img_url varchar(150),
                        mail varchar(150) not null unique,
                        password varchar(255) not null
);

CREATE TABLE "event" (
                         id serial not null unique,
                         title varchar(255) not null,
                         description varchar(500) not null,
                         text varchar(1000) not null,
                         city varchar(255) not null,
                         category varchar(255) not null,
                         viewed BIGINT not null,
                         img_url varchar(500),
                         date varchar(10) not null,
                         geo varchar(255) not null,
                         tag varchar(30)[],
                         author_id int references "user" (id) on delete cascade not null
);

CREATE TABLE "view" (
                        id serial not null unique,
                        event_id int references "event" (id) on delete cascade not null,
                        user_id int references "user" (id) on delete cascade,
                        view_time timestamp
);

CREATE TABLE "visitor" (
                        id serial not null unique,
                        event_id int references "event" (id) on delete cascade not null,
                        user_id int references "user" (id) on delete cascade,
                        agree_to_go_date date
);

INSERT INTO "user" (name, surname, mail, password, about)
VALUES ('Andrey', 'Ivanov', 'test@mail.ru', 'hashhashhash', 'I am soooo coool DAAAMN'),
       ('Ivan', 'Andreev', 'kool@mail.ru', 'hashhashhash', 'My name is Ivan Andreev'),
       ('Petr', 'Leshiy', 'clown@mail.ru', 'hashhashhash', 'Yo my desctiption');

INSERT INTO "event" (title, description, text, city, category, viewed, img_url, date, geo, tag, author_id)
VALUES ('Tusovka', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Концерты', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array['Stupid', 'Alcohol'], 1),
       ('Funny party', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Вечеринки', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array['Alcohol', 'Hype'], 3),
       ('Really funny party', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Тусовки', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array ['Hype', 'Boomerang'], 2),
       ('PAPAPAPA', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Выставки', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array['Boring', 'Alcohol'], 1),
       ('Tusovka', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Театры', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array['Stupid', 'Alcohol'], 1),
       ('Funny party', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array['Alcohol', 'Hype'], 3),
       ('Really funny party', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array ['Hype', 'Boomerang'], 2),
       ('PAPAPAPA', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array['Boring', 'Alcohol'], 1),
       ('Tusovka', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array['Stupid', 'Alcohol'], 1),
       ('Funny party', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array['Alcohol', 'Hype'], 3),
       ('Really funny party', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array ['Hype', 'Boomerang'], 2),
       ('PAPAPAPA', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array['Boring', 'Alcohol'], 1),
       ('Tusovka', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array['Stupid', 'Alcohol'], 1),
       ('Funny party', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array['Alcohol', 'Hype'], 3),
       ('Really funny party', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array ['Hype', 'Boomerang'], 2),
       ('PAPAPAPA', 'Funny party this time', 'YOOOOOOOOOOOOOOOOOOOOOOOOOOOo', 'Moscow', 'Party', 123, 'img.png', '10.10.2008', 'ul Pushkina',
        array['Boring', 'Alcohol'], 1);