DROP VIEW IF EXISTS users_view;
CREATE VIEW users_view AS
    SELECT 
        u.id
        , u.name
        , u.email
        , u.created_at
        , u.updated_at
        , u.deleted_at
    FROM users u;

    DROP VIEW IF EXISTS orders_view;
    CREATE VIEW orders_view AS
        SELECT
            o.id,
            o.user_id,
            o.status,
            u.name AS user_name,
            u.email AS user_email,
            o.created_at,
            o.updated_at,
            o.deleted_at
        FROM orders o
        LEFT JOIN users u ON o.user_id = u.id
        WHERE o.deleted_at IS NULL;
                

    DROP VIEW IF EXISTS order_items_view;
    CREATE VIEW order_items_view AS
        SELECT
            oi.id,
            oi.order_id,
            oi.product_id,
            oi.quantity AS ordered_quantity,
            oi.price AS items_price,
            p.price AS product_price,
            p.name AS product_name,
            p.stock AS remaining_product_stock
        FROM order_items oi
        LEFT JOIN products p ON oi.product_id = p.id;