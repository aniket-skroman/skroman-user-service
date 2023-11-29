create table "skroman_client" (
    "id" uuid default uuid_generate_v4 () primary key,
    "user_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL ,
    "password" varchar,
    "contact" varchar NOT NULL,
    "address" varchar NOT NULL,
    "city" varchar,
    "state" varchar,
    "pincode" varchar,
    "created_at" timestamptz NOT NULL default (now()),
    "updated_at" timestamptz NOT NULL default (now())
)