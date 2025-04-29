create table if not exists public.tasks
(
    id BIGSERIAL primary key,
    name text not null,
    description text,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

