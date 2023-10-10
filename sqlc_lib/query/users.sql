-- name: CreateNewUser :one
insert into users (
    full_name,
    email,
    password,
    contact,
    user_type
) values (
    $1,$2,$3,$4,$5
) returning *;


-- name: GetUserByEmailOrContact :one
select * from users
where email=$1 or contact = $1
limit 1;