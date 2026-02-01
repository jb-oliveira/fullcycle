-- name: ListCategories :many
-- name: ListCategoriesDynamic :many
SELECT id, name, description
FROM categories
ORDER BY
    -- Handle Ascending Cases
    CASE WHEN @sort_order::text <> 'DESC' THEN
        CASE 
            WHEN @sort_column::text = 'id' THEN id
            WHEN @sort_column::text = 'name' THEN name
            WHEN @sort_column::text = 'description' THEN description
        END
    END ASC,
    -- Handle Descending Cases
    CASE WHEN @sort_order::text = 'DESC' THEN
        CASE 
            WHEN @sort_column::text = 'id' THEN id
            WHEN @sort_column::text = 'name' THEN name
            WHEN @sort_column::text = 'description' THEN description
        END
    END DESC
LIMIT @limit_val
OFFSET @offset_val;

-- name: GetCategory :one
SELECT id, name, description
FROM categories
WHERE id = $1;

-- name: CreateCategory :one
INSERT INTO categories (id, name, description)
VALUES ($1, $2, $3)
RETURNING id, name, description;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, description = $3
WHERE id = $1
RETURNING id, name, description;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;

-- name: CountCategories :one
SELECT count(*) FROM categories;


-- name: CreateCategoriesBatch :copyfrom
INSERT INTO categories (id, name, description)
VALUES ($1, $2, $3);

-- name: CreateCourse :one
insert into courses (id, category_id, name, description, price)
values ($1,$2,$3,$4,$5)
RETURNING id, category_id, name, description, price;