CREATE TABLE "user" (
                        id serial not null unique,
                        name varchar(50) not null,
                        surname varchar(50) not null,
                        about varchar(150) default '',
                        img_url varchar(150) default '',
                        mail varchar(150) not null unique,
                        password varchar(255) not null
);

CREATE TABLE "event" (
                         id serial not null unique,
                         title varchar(255) not null,
                         description varchar(500) not null,
                         text varchar(2200) not null,
                         city varchar(255) not null,
                         category varchar(255) not null,
                         viewed BIGINT not null,
                         img_url varchar(500) default '',
                         date varchar(10) not null,
                         geo varchar(255) default '' not null,
                         address varchar(255) default '' not null,
                         tag varchar(30)[],
                         author_id int references "user" (id) on delete cascade not null
);

CREATE TABLE "view" (
                        id serial not null unique,
                        event_id int references "event" (id) on delete cascade not null,
                        user_id int references "user" (id) on delete cascade,
                        UNIQUE(event_id, user_id),
                        date date
);

CREATE TABLE "visitor" (
                        id serial not null unique,
                        event_id int references "event" (id) on delete cascade not null,
                        user_id int references "user" (id) on delete cascade not null,
                        UNIQUE(event_id, user_id),
                        date date
);

CREATE TABLE "subscribe" (
                           id serial not null unique,
                           subscribed_id int references "user" (id) on delete cascade not null,
                           subscriber_id int references "user" (id) on delete cascade not null,
                           UNIQUE(subscribed_id, subscriber_id),
                           CHECK ( subscribed_id <> subscribe.subscriber_id )
);

drop table "notification";

CREATE TABLE "notification" (
    id serial not null unique,
    type varchar(50) CHECK (type in ('0', '1', '2', '3')) not null,
                                    receiver_id varchar(50) not null,
                                    user_id varchar(50) not null,
                                    user_name varchar(50) not null,
                                    user_surname varchar(50) not null,
                                    user_img_url varchar(150) not null,
                                    event_id varchar(50) default '',
                                    event_title varchar(255) default '',
                                    seen bool default false,
                                    UNIQUE(type, receiver_id, user_id, event_id)
);