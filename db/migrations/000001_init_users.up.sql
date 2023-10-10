create table users (
    "id" uuid default uuid_generate_v4 () primary key,
    "full_name" varchar UNIQUE NOT NULL,
    "email" varchar UNIQUE NOT NULL ,
    "password" varchar NOT NULL,
    "contact" varchar NOT NULL,
    "user_type" varchar NOT NULL default 'guest_user',
    "created_at" timestamptz NOT NULL default (now()),
    "updated_at" timestamptz NOT NULL default (now())
);

alter table users
    add constraint check_full_name check (full_name <> ''),
    add constraint check_email check (email <> '' and email LIKE '%@%.%' AND email NOT LIKE '@%' AND email NOT LIKE '%@%@%'),
    add constraint check_password check (password <> '');