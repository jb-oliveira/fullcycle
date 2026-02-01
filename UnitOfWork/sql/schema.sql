create table public.categories
(
    id          varchar(36)  not null
        primary key,
    name        varchar(255) not null
        unique,
    description text
);

alter table public.categories
    owner to postgres;

create table public.courses
(
    id          varchar(36)  not null
        primary key,
    category_id varchar(36)  not null
        references public.categories,
    name        varchar(255) not null
        unique,
    description text,
    price       numeric      not null
);

alter table public.courses
    owner to postgres;

