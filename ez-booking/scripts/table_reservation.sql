create table roles
(
    id   serial       not null
        constraint roles_pkey
            primary key,
    name varchar(255) not null
);

insert into roles (name)
values ('admin'),
       ('user');


create table users
(
    id       serial       not null
        constraint users_pkey
            primary key,
    name     varchar(255) not null,
    role_id  int          not null,
    email    varchar(255) not null
        constraint users_uc_email
            unique,
    mobile   varchar(255) not null,
    password char(60)     not null,
    created  timestamp    not null,
    CONSTRAINT fk_role_id
        FOREIGN KEY (role_id)
            REFERENCES roles (id)
);

alter table users add column profile_image_url varchar(255) default '/static/img/default_profile_image.png';

create index users_name_password_idx
    on users (email, password);


create table locations
(
    id   serial       not null primary key,
    name varchar(255) not null
);

insert into locations (name)
values ('default'),
       ('sofa'),
       ('outside'),
       ('bar');


create table cafes
(
    id   serial       not null primary key,
    name varchar(255) not null
);

insert into cafes (name)
values ('tasty_food'),
       ('random_cafe');

create table tables
(
    id          int not null,
    cafe_id     int not null,
    capacity    int not null,
    location_id int not null,
    CONSTRAINT ck_id_cafe_id UNIQUE (id, cafe_id),
    constraint fk_location_id
        foreign key (location_id)
            references locations (id)
);

alter table tables
    add constraint tables_fk_cafe_id
        foreign key (cafe_id)
            references cafes (id);


insert into tables (id, cafe_id, capacity, location_id)
values (1, 1, 2, 4),
       (2, 1, 2, 4),
       (3, 1, 4, 4),
       (4, 1, 4, 2),
       (5, 1, 4, 2),
       (6, 1, 8, 2),
       (7, 1, 8, 2),
       (8, 1, 4, 3),
       (9, 1, 4, 3),
       (10, 1, 2, 3),
       (11, 1, 2, 3),
       (12, 1, 2, 1),
       (13, 1, 2, 1),
       (14, 1, 4, 1),
       (15, 1, 4, 1);

create table events
(
    id   serial       not null primary key,
    name varchar(255) not null
);

insert into events (name)
values ('default'),
       ('romantic dinner'),
       ('birthday');


create table reservations
(
    id                serial       not null primary key,
    user_id           int,
    cafe_id           int          not null,
    table_id          int          not null,
    cust_name         varchar(255) not null,
    cust_mobile       varchar(255) not null,
    cust_email        varchar(255) not null,
    event_id          int          not null default 1,
    event_description text,
    num_of_persons    int          not null,
    date              timestamp    not null,
    CONSTRAINT reservations_fk_cafe_id_table_id
        FOREIGN KEY (cafe_id, table_id)
            REFERENCES tables (cafe_id, id),
    constraint reservations_fk_user_id
        foreign key (user_id)
            references users (id),
    constraint fk_event_id
        foreign key (event_id)
            references events (id)
);

create table cafes_events
(
    cafe_id  int not null,
    event_id int not null,
    constraint cafes_events_fk_cafe_id foreign key (cafe_id) references cafes (id),
    constraint cafes_events_fk_event_id foreign key (event_id) references events (id),
    CONSTRAINT ck_cafe_id_event_id primary key (cafe_id, event_id)
);

insert into cafes_events (cafe_id, event_id)
VALUES (1, 1),
       (1, 2),
       (1, 3);

