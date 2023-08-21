
Create database BlogDB;
create schema public;
CREATE TABLE public.blog_posts (
    id serial PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL
);



postgres