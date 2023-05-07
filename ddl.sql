create table "user"
(
    id         serial
        primary key,
    user_id    bigint                  not null
        unique,
    username   varchar(255),
    firstname  varchar(255),
    lastname   varchar(255),
    warn_count integer   default 0     not null,
    created    timestamp default now() not null,
    modified   timestamp default now() not null,
    enabled    boolean   default true  not null,
    author     varchar(255)
);

comment on table "user" is 'Таблица с пользователями в телеграм';

alter table "user"
    owner to gonec;

create table moderator
(
    id              serial
        primary key,
    name            varchar(255)                    not null
        unique,
    description     text,
    pattern         text,
    words           text,
    delete          boolean   default true          not null,
    warn            boolean   default true          not null,
    warn_number     integer   default '-1'::integer not null,
    enabled         boolean   default true          not null,
    created         timestamp default now()         not null,
    modified        timestamp default now()         not null,
    author          varchar(255),
    until_date      bigint    default 86400         not null,
    revoke_messages boolean   default false         not null
);

alter table moderator
    owner to gonec;

create table chat
(
    id          serial
        primary key,
    chat_id     bigint                  not null
        unique,
    title       varchar(255),
    type        varchar(255),
    description text,
    created     timestamp default now() not null,
    modified    timestamp default now() not null,
    author      varchar(255),
    enabled     boolean   default true  not null
);

alter table chat
    owner to gonec;

create table bot
(
    id           serial
        primary key,
    name         varchar(255)                                      not null
        unique,
    token        varchar(1024)                                     not null
        unique,
    description  text,
    created      timestamp   default now()                         not null,
    modified     timestamp   default now()                         not null,
    author       varchar(255),
    enabled      boolean     default true                          not null,
    timeout      integer     default 60                            not null,
    parse_mode   varchar(10) default 'Markdown'::character varying not null,
    welcome      boolean     default false                         not null,
    welcome_text text
);

alter table bot
    owner to gonec;

create table bot_moderator
(
    bot_id       integer not null
        constraint bot_moderator_bot_id_fk
            references bot
            on update cascade on delete cascade,
    moderator_id integer not null
        constraint bot_moderator_moderator_id_fk
            references moderator
            on update cascade on delete cascade
);

alter table bot_moderator
    owner to gonec;

create unique index bot_moderator_bot_id_moderator_id_uindex
    on bot_moderator (bot_id, moderator_id);

create view vw_bot_moderator
            (bot_id, bot, token, timeout, parse_mode, welcome, welcome_text, moderator_id, name, pattern, delete, warn,
             warn_number, words)
as
SELECT bm.bot_id,
       b.name AS bot,
       b.token,
       b.timeout,
       b.parse_mode,
       b.welcome,
       b.welcome_text,
       bm.moderator_id,
       m.name,
       m.pattern,
       m.delete,
       m.warn,
       m.warn_number,
       m.words
FROM telegram.bot_moderator bm
         JOIN telegram.bot b ON bm.bot_id = b.id
         JOIN telegram.moderator m ON bm.moderator_id = m.id;

alter table vw_bot_moderator
    owner to gonec;

