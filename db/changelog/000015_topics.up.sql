create table if not exists public.topics (
                                             id bigserial primary key,
                                             topic_title text not null,
                                             topic_content text not null,
                                             preview_url text,
                                             created_at timestamp without time zone,
                                             updated_at timestamp without time zone
)