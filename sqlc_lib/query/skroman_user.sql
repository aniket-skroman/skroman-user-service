/* create/register new user */
-- name: CreateSkromanUser :one
insert into skroman_client (
    user_name,
    email,
    password,
    contact,
    address,
    city,
    state,
    pincode
) values (
    $1,$2,$3,$4,$5,$6,$7,$8
) returning *;


-- name: FetchAllClients :many
select * from skroman_client
order by created_at desc    
limit $1
offset $2;

-- name: CountOFClients :one
select count(*) from skroman_client;

-- name: DeleteClient :execresult
delete from skroman_client
where id = $1;

-- name: FetchClientById :one
select * from skroman_client
where id = $1
limit 1;

/* update skroman client info */
-- name: UpdateSkromanClientInfo :one
update skroman_client
set user_name=$2,
email=$3,password=$4,contact=$5,
address=$6,city=$7,state=$8,pincode=$9,
updated_at = CURRENT_TIMESTAMP
where id=$1
returning * ;