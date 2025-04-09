create table if not exists public.collection (
    id bigserial primary key,
    category_id bigint,
    name text not null,
    description text not null,
    background_image text not null,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

create table if not exists public.collection_item (
    id bigserial primary key,
    collection_id bigint not null REFERENCES public.collection (id),
    item_id bigint not null REFERENCES public.item (id)
);

alter table public.item add column popular integer;
alter table public.item add column new integer;
alter table public.item add column unavailable boolean not null default false;