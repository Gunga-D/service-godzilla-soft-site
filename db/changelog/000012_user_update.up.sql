ALTER TABLE public.user ALTER COLUMN email DROP NOT NULL;
ALTER TABLE public.user ALTER COLUMN password DROP NOT NULL;
ALTER TABLE public.user ADD COLUMN photo_url text NULL; 
ALTER TABLE public.user ADD COLUMN username text NULL;
ALTER TABLE public.user ADD COLUMN first_name text NULL;
ALTER TABLE public.user ADD COLUMN telegram_id bigint NULL;
ALTER TABLE public.user ADD COLUMN steam_link text NULL;