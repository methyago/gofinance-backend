-- name: CreateAccount :one
INSERT INTO accounts (
    user_id,
    category_id,
    title,
    type,
    description,
    date,
    value
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts WHERE id = $1 LIMIT 1;

-- name: GetAccounts :many
SELECT a.id, a.user_id, 
       a.title, a.type, a.description, 
       a.value, a.date, a.created_at, 
       c.title as category_title
  FROM accounts a
  LEFT JOIN categories c on c.id = a.category_id
 WHERE a.user_id = @user_id 
   AND a.type = @type
   AND a.category_id = COALESCE(sqlc.narg('category_id'), a.category_id)
   AND (UPPER(a.title) LIKE CONCAT('%', UPPER(@title::text), '%'))
   AND (UPPER(a.description) LIKE CONCAT('%', UPPER(@description::text), '%'))
   AND a.date = COALESCE(sqlc.narg('date'), a.date);

-- name: GetAccountsReports :one
SELECT SUM(value) AS sum_value FROM accounts 
WHERE user_id = $1 AND type = $2;

-- name: GetAccountGraph :one
SELECT COUNT(*) FROM accounts
WHERE user_id = $1 and type = $2;

-- name: UpdateAccounts :one
UPDATE accounts SET title = $2, description = $3, value = $4 WHERE id = $1 RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;