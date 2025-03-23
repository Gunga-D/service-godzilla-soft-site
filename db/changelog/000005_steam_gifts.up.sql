alter table public.item add column is_steam_gift boolean default false;
alter table public.item add column steam_app_id bigint;