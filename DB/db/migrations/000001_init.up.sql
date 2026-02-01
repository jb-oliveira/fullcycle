create table categories (
    id varchar(36) primary key,
    name varchar(255) not null,
    description text
);

create table courses (
    id varchar(36) primary key,
    category_id varchar(36) not null references categories(id),
    name varchar(255) not null,
    description text,
    price decimal not null    
);