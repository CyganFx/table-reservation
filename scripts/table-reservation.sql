create table roles
(
    id   serial       not null
        constraint roles_pkey
            primary key,
    name varchar(255) not null
);

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
