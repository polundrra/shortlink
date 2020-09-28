create table if not exists link (
    id serial primary key,
    url varchar(2048) not null,
    code varchar(32) unique,
    is_custom boolean not null
);

create sequence if not exists seq;
