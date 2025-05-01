create table if not exists tasks
(
    id BIGSERIAL primary key,
    name text not null,
    description text
);