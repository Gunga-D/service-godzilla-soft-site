create table if not exists public.finished_neuro_task (
    id uuid primary key,
    query text not null,
    result text not null,
    created_at timestamp without time zone
);