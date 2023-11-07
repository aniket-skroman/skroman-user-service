alter table users
    add column department varchar not null default 'skroman';

alter table users
    add constraint check_department check (department <> '');