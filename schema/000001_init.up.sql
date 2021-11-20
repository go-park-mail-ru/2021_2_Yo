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

/*
Просмотр (<=> избранное)
event_id - id мероприятия
user_id - id посетителя
date - для графиков
*/
CREATE TABLE "view" (
                        id serial not null unique,
                        event_id int references "event" (id) on delete cascade not null,
                        user_id int references "user" (id) on delete cascade,
                        UNIQUE(event_id, user_id),
                        date date
);

/*
Посетитель (<=> избранное)
event_id - id мероприятия
user_id - id посетителя
date - для графиков
*/
CREATE TABLE "visitor" (
                        id serial not null unique,
                        event_id int references "event" (id) on delete cascade not null,
                        user_id int references "user" (id) on delete cascade not null,
                        UNIQUE(event_id, user_id),
                        date date
);

/*
Подписка
subscriber_id - id пользователя, который подписался
subscribed_id - id пользователя, на которого подписались
*/
CREATE TABLE "subscribe" (
                           id serial not null unique,
                           subscribed_id int references "user" (id) on delete cascade not null,
                           subscriber_id int references "user" (id) on delete cascade not null,
                           UNIQUE(subscribed_id, subscriber_id),
                           CHECK ( subscribed_id <> subscribe.subscriber_id )
);