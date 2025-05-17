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