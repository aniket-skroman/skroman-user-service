create table user_fcm_data(
    "id" uuid default uuid_generate_v4 () primary key,
    "user_id" uuid not null,
    "fcm_token" varchar not null,
    "created_at" timestamptz NOT NULL default (now()),
    "updated_at" timestamptz NOT NULL default (now()),
    foreign key ("user_id") references users("id")
);