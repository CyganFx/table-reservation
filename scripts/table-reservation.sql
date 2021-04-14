create table users
(
    id       serial       not null
        constraint users_pkey
            primary key,
    name     varchar(255) not null,
    email    varchar(255) not null
        constraint users_uc_email
            unique,
    mobile   varchar(255) not null,
    password char(60)     not null,
    created  timestamp    not null
);

alter table users
    owner to postgres;

create index users_name_password_idx
    on users (email, password);