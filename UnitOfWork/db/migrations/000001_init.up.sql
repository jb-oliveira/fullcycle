create table if not exists categories (
    id varchar(36) primary key,
    name varchar(255) not null unique
);

create table if not exists courses (
    id varchar(36) primary key,
    category_id varchar(36) not null references categories(id),
    name varchar(255) not null unique,
    description text,
    price decimal not null    
);