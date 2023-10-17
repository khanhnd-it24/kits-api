create table if not exists users (
    "id" serial primary key,
    "username" text not null,
    "password" text not null,
    "full_name" text,
    "created_at" timestamp not null,
    "updated_at" timestamp not null,
    "deleted_at" timestamp
);

create table if not exists refresh_tokens (
    "id" serial primary key,
    "user_id" integer not null,
    "token" text not null,
    "expired_at" timestamp not null,
    "created_at" timestamp not null,
    "updated_at" timestamp not null
);

create table if not exists habits (
    "id" serial primary key,
    "name" text not null,
    "order" integer not null,
    "icon" text not null,
    "theme_color" text not null,
    "user_id" integer not null,
    "created_at" timestamp not null,
    "updated_at" timestamp not null
);

create table if not exists daily_habits (
    "habit_id" integer,
    "date" text,
    primary key ("habit_id", "date")
);

create table if not exists passwords (
    "id" serial primary key,
    "name" text not null,
    "password" text not null,
    "notes" text,
    "user_id" integer not null,
    "created_at" timestamp not null,
    "updated_at" timestamp not null
);

create index on daily_habits("habit_id");
alter table habits add foreign key ("user_id") references users("id");
alter table refresh_tokens add foreign key ("user_id") references users("id");
alter table passwords add foreign key ("user_id") references users("id");
alter table daily_habits add foreign key ("habit_id") references habits("id");