alter table users
    add column emp_code varchar not null unique default substring(md5(random()::text), 0, 8);