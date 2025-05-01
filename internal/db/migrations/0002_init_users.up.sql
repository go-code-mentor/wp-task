create table if not exists users
(
    id BIGSERIAL primary key,
    login varchar(64) not null unique
);