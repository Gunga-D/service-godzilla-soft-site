create type e_roulette_session_status as enum
    ('wait-for-payment', 'ready-to-roll', 'out-of-games');

create type e_random_item_catetory as enum
   ('common', 'uncommon', 'rare', 'special', 'golden');

-- stores informtaion about user roulette session
create table roulette_session
(
    id         UUID PRIMARY KEY          DEFAULT gen_random_uuid(),
    status     e_roulette_session_status default 'wait-for-payment' not null,
    created_at timestamp not null        default current_timestamp,
    updated_at timestamp not null        default current_timestamp
);


-- stores references to items from catalog as well as additional cost and category
create table roulette_random_item
(
    id            serial primary key     not null references public.item (id),
    total_cost    int                    not null,
    item_category e_random_item_catetory not null
);

-- stores items related to specific session as well as their states
create table roulette_session_items
(
    session_id UUID REFERENCES roulette_session (id) ON DELETE CASCADE,
    item_id    int       not null references roulette_random_item (id),
    is_top     bool      not null default false,
    won_at     timestamp null,
    constraint pk_session_items unique (session_id, item_id)
);