CREATE TABLE users (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT  CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE user_sessions (
    token CHAR(64) PRIMARY KEY,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE orders (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL,
    Status VARCHAR(20) NOT NULL,
    shipping_address VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT  CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX ON orders (deleted_at) WHERE deleted_at IS NULL;

CREATE TABLE products (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT  CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX ON products (deleted_at) WHERE deleted_at IS NULL;

CREATE TABLE order_items (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) 
);

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
            o.shipping_address,
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