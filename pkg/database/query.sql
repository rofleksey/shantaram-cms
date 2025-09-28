-- name: CreateOrder :one
INSERT INTO orders (id, table_id, client_name, client_phone, client_comment, status, seen, items)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetOrderByID :one
SELECT *
FROM orders
WHERE id = $1;

-- name: GetOrdersPaginated :many
SELECT *
FROM orders
ORDER BY index DESC
OFFSET $1 LIMIT $2;

-- name: CountOrders :one
SELECT COUNT(*)
FROM orders;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status  = $2,
    updated = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: SetOrderSeen :exec
UPDATE orders
SET seen    = true,
    updated = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteOrder :exec
DELETE
FROM orders
WHERE id = $1;

-- name: GetMenus :many
SELECT *
FROM menu
ORDER BY created;

-- name: CreateProductGroup :exec
INSERT INTO product_groups (id, menu_id, title, index)
VALUES (@id, @menu_id::VARCHAR(255), @title,
        (SELECT COALESCE(MAX(index), 0) + 1 FROM product_groups WHERE menu_id = @menu_id:: VARCHAR (255)) );

-- name: GetProductGroupByID :one
SELECT *
FROM product_groups
WHERE id = $1;

-- name: GetAllProductGroups :many
SELECT *
FROM product_groups
ORDER BY index;

-- name: UpdateProductGroup :exec
UPDATE product_groups
SET title   = $2,
    updated = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateProductGroupIndex :exec
UPDATE product_groups
SET index   = $2,
    updated = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteProductGroup :exec
DELETE
FROM product_groups
WHERE id = $1;

-- name: CreateProduct :exec
INSERT INTO products (id, group_id, title, description, price, index)
VALUES (@id, @group_id::UUID, @title, @description, @price,
        (SELECT COALESCE(MAX(index), 0) + 1 FROM products WHERE group_id = @group_id::UUID) );

-- name: GetProductByID :one
SELECT *
FROM products
WHERE id = $1;

-- name: GetProductsByGroup :many
SELECT *
FROM products
WHERE group_id = $1
ORDER BY index;

-- name: GetAllProducts :many
SELECT *
FROM products
ORDER BY available DESC, index, group_id;

-- name: UpdateProduct :exec
UPDATE products
SET title       = $2,
    description = $3,
    price       = $4,
    available   = $5,
    updated     = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateProductIndex :exec
UPDATE products
SET index   = $2,
    updated = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE
FROM products
WHERE id = $1;

-- name: SetProductAvailability :exec
UPDATE products
SET available = $2,
    updated   = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: SearchProducts :many
SELECT *
FROM products
WHERE (title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
  AND available = true
ORDER BY title;


-- name: GetMigrations :many
SELECT *
FROM migration
ORDER BY id;

-- name: CreateMigration :one
INSERT INTO migration (id, applied)
VALUES ($1, $2) RETURNING id;
