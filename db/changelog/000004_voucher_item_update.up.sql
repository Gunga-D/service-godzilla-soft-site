alter table public.item add column bx_image_url text;
alter table public.item add column limit_price bigint;
alter table public.item add column bx_gallery_urls text[];

create table if not exists public.voucher (
    id bigserial primary key,
    type text not null,
    value text not null,
    impact bigint not null,
    has_activated boolean not null,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

alter table public.item add column yandex_id text;
