-- name: CreateNewUser :one
insert into users (
    full_name,
    email,
    password,
    contact,
    user_type,
    department
) values (
    $1,$2,$3,$4,$5,$6
) returning *;


-- name: GetUserByEmailOrContact :one
select * from users
where email=$1 or contact = $1
limit 1;

-- name: CheckEmailOrContactExists :execrows
select * from users
where email=$1 or contact=$2
limit 1;

-- name: CountUsers :one
select count(*) from users;

-- name: FetchAllUsers :many
select * from users 
order by created_at desc 
limit $1
offset $2;

-- name: UpdateUser :execresult
update users 
set full_name=$2,
contact=$3,
user_type=$4,
updated_at = CURRENT_TIMESTAMP
where id=$1 
returning *;

-- name: CheckForContact :one
select * from users 
where contact = $1
and id <> $2;

-- name: DeleteUser :execrows
delete from users 
where id = $1;


-- name: GetUserById :one
select * from users
where id = $1
limit 1;


/* fetch users by department */
-- name: UsersByDepartment :many
select * from users 
where department = $1
group by id
order by created_at desc 
limit $2
offset $3;

/* count users by department */
-- name: CountUsersByDepartment :one
select count(*) from users 
where department = $1;