-- name: CreateCategory :exec
insert into categories (id, name, description) values ($1, $2, $3);

-- name: CreateCourse :exec
insert into courses (id, category_id, name, description, price) 
values ($1, $2, $3, $4, $5);

-- name: FindCategoryByID :one
select * from categories where id = $1;

-- name: FindCourseByID :one
select * from courses where id = $1;
