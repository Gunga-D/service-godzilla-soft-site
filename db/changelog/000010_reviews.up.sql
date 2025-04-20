create table if not exists public.review (
    id bigserial primary key,
    user_id bigint,
    item_id bigint not null,
    comment text,
    score int not null,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);