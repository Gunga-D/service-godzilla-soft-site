CREATE EXTENSION IF NOT EXISTS pgcrypto;

create table if not exists public.order (
    id uuid primary key default gen_random_uuid(),
    email text not null,
    code_value text not null,
    amount  bigint not null,
    status text not null,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

create table if not exists public.user (
    id bigserial primary key,
    email text not null,
    password text not null,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

create table if not exists public.item (
    id bigserial primary key,
    title text not null,
    description text,
    category_id bigint,
    platform text not null,
    region text not null,
    current_price bigint not null,
    is_for_sale boolean not null,
    old_price bigint,
    thumbnail_url text not null,
    status text not null,
    slip text not null,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

create table if not exists public.code (
	value text primary key,
    item_id bigint not null,
    status text not null,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);