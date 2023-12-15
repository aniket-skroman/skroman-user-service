-- name: CreateUserFCMData :one
insert into user_fcm_data (
    user_id,
    fcm_token
) values (
    $1,$2
) returning *;

-- name: FetchFCMTokensByUser :many
select * from user_fcm_data
where user_id = $1;

