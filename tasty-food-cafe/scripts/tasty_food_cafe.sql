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


create table tables
(
    id          serial not null primary key,
    capacity    int    not null,
    location_id int    not null,
    constraint fk_location_id
        foreign key (location_id)
            references locations (id)
);

insert into tables (capacity, location_id)
values (2, 4),
       (2, 4),
       (4, 4),
       (4, 2),
       (4, 2),
       (8, 2),
       (8, 2),
       (4, 3),
       (4, 3),
       (2, 3),
       (2, 3),
       (2, 1),
       (2, 1),
       (4, 1),
       (4, 1);


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
    id             serial       not null primary key,
    table_id       int          not null,
    cust_name      varchar(255) not null,
    cust_mobile    varchar(255) not null,
    cust_email     varchar(255) not null,
    event_id       int          not null default 1,
    num_of_persons int          not null,
    date           timestamp    not null,
    CONSTRAINT fk_table_id
        FOREIGN KEY (table_id)
            REFERENCES tables (id),
    constraint fk_event_id
        foreign key (event_id)
            references events (id)
);