create table if not exists link (
    id serial primary key,
    url varchar(2048) not null unique
    code varchar(32) not null unique
);

create sequence if not exists seq;
