begin;

create table if not exists namespaces (
    id text not null,
    token text not null,
    active bool not null default true,
    created bigint not null default date_part('epoch'::text, now()),
    primary key (id)
);

create type enum_element_type as enum (
    'PASTE',
    'REDIRECT'
);

create table if not exists elements (
    namespace text not null,
    key text not null,
    type enum_element_type not null,
    internal_data jsonb not null,
    public_data jsonb not null,
    views int not null default 0,
    max_views int not null default -1,
    valid_from bigint not null default -1,
    valid_until bigint not null default -1,
    created bigint not null default date_part('epoch'::text, now()),
    primary key (namespace, key)
);

create table if not exists invites (
    code text not null,
    uses int not null default 0,
    max_uses int not null default -1,
    created bigint not null default date_part('epoch'::text, now()),
    primary key (code)
);

commit;
