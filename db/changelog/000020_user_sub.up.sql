create table if not exists public.subscription_bills (
    id uuid primary key default gen_random_uuid(),
    user_id bigint not null,
    amount  bigint not null,
    rebill_id text,
    status text not null,
    need_prolong boolean not null default true,
    created_at bigint not null,
    updated_at timestamp without time zone,
    expired_at bigint not null
);

alter table public.item add column in_sub boolean not null default false;

create table if not exists public.accounts (
    id bigserial primary key,
    item_id bigint not null,
    login text not null,
    password text not null
);
