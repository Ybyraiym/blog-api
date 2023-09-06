
Create database BlogDB;
create schema public;
CREATE TABLE public.blog_posts (
    id serial PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL
);

CREATE TABLE public.users (
    id serial PRIMARY KEY,
    Login VARCHAR(255) Unique NOT NULL,
    Username VARCHAR(255) Unique NOT NULL,
    Password  VARCHAR(255) NOT NULL
);




postgres