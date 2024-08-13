set statement_timeout=0;
set lock_timeout=0;
set idle_in_transaction_session_timeout=0;
set client_encoding='UTF-8';
set standard_conforming_strings=on;
set client_min_messages=warning;
set row_security=off;

create extension if not exists plpgsql with schema pg_catalog;
create extension if not exists "uuid-ossp" with schema pg_catalog;

set search_path=public, pg_catalog;
set default_tablespace='';


--users

create extension if not exists pgcrypto;

create table users(
    id uuid not null default uuid_generate_v1mc(),
    username text not null unique,
    user_password text not null,
    user_role text not null,
    access_token text,
    constraint users_pk primary key (id) 
);

create index user_access_token
on users (access_token);

insert into users(username, user_password,user_role)
values
('gurpreet',crypt('gurpreet',gen_salt('bf')),'admin'),
('aastha',crypt('aastha',gen_salt('bf')),'user'),
('harshit',crypt('harshit',gen_salt('bf')),'user');

--chats

CREATE TABLE chats (
    id uuid not null default uuid_generate_v1mc(),
    message TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    username TEXT NOT NULL,
    CONSTRAINT fk_username FOREIGN KEY (username) REFERENCES users (username)
);
create index chats_username
on chats (username);