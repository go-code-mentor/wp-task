create table if not exists access_tokens
(
    id BIGSERIAL primary key,
    user_id BIGSERIAL references users(id) not null,
    token varchar(128) unique not null
);