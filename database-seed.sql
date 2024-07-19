-- Create users table
CREATE TABLE IF NOT EXISTS public.users
(
    id                serial primary key       not null,
    created_at        timestamp with time zone default current_timestamp,
    updated_at        timestamp with time zone default current_timestamp,
    deleted_at        timestamp with time zone,
    username          text                     not null,
    email             text                     not null,
    password          text                     not null,
    registration_date timestamp with time zone not null,
    last_login        timestamp with time zone,
    avatar            text
);


-- Create indexes for users table
CREATE UNIQUE INDEX IF NOT EXISTS uni_users_email ON public.users USING btree (email);
CREATE UNIQUE INDEX IF NOT EXISTS uni_users_username ON public.users USING btree (username);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON public.users USING btree (deleted_at);

-- Create profiles table
CREATE TABLE IF NOT EXISTS public.profiles
(
    id         serial primary key not null,
    deleted_at timestamp with time zone,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    first_name text,
    last_name  text,
    headline   text,
    language   text,
    website    text,
    twitter    text,
    facebook   text,
    youtube    text,
    github     text,
    about_me   text,
    user_id    bigint             not null,
    foreign key (user_id) references public.users (id)
);

-- Create indexes for profiles table
CREATE UNIQUE INDEX IF NOT EXISTS uni_profiles_user_id ON public.profiles USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_profiles_user_id ON public.profiles USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_profiles_deleted_at ON public.profiles USING btree (deleted_at);

-- Create categories table
CREATE TABLE IF NOT EXISTS public.categories
(
    id          serial primary key not null,
    created_at  timestamp with time zone default current_timestamp,
    updated_at  timestamp with time zone default current_timestamp,
    deleted_at  timestamp with time zone,
    title       varchar            not null,
    slug        varchar            not null,
    image       varchar,
    description varchar,
    type        varchar                  default 'category'
);


-- Create indexes for categories table
CREATE UNIQUE INDEX IF NOT EXISTS uni_categories_title ON public.categories USING btree (title);
CREATE UNIQUE INDEX IF NOT EXISTS uni_categories_slug ON public.categories USING btree (slug);
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON public.categories USING btree (deleted_at);


DROP TABLE IF EXISTS public.courses;
-- Create courses table
CREATE TABLE IF NOT EXISTS public.courses
(
    id           serial primary key not null,
    course_id    uuid               not null,
    created_at   timestamp with time zone default current_timestamp,
    updated_at   timestamp with time zone default current_timestamp,
    deleted_at   timestamp with time zone,
    thumbnail    text,
    title        text               not null unique,
    slug         text               not null,
    description  text               not null,
    publish_date timestamp with time zone,
    price        double precision   not null
);

-- Create indexes for courses table
CREATE UNIQUE INDEX IF NOT EXISTS uni_courses_slug ON public.courses USING btree (slug);
CREATE INDEX IF NOT EXISTS idx_courses_deleted_at ON public.courses USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_courses_course_id ON public.courses USING btree (course_id);

-- Create topics table
CREATE TABLE IF NOT EXISTS public.topics
(
    id         serial primary key not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone,
    title      text               not null,
    slug       text               not null
);

-- Create indexes for topics table
CREATE UNIQUE INDEX IF NOT EXISTS uni_topics_title ON public.topics USING btree (title);
CREATE UNIQUE INDEX IF NOT EXISTS uni_topics_slug ON public.topics USING btree (slug);
CREATE INDEX IF NOT EXISTS idx_topics_deleted_at ON public.topics USING btree (deleted_at);



create table authors_courses
(
    id        serial primary key,
    author_id int  not null,
    course_id uuid not null,
    CONSTRAINT  unique_author_course unique (author_id, course_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS uni_authors_courses ON public.authors_courses USING btree (author_id);
CREATE UNIQUE INDEX IF NOT EXISTS uni_authors_courses ON public.authors_courses USING btree (course_id);




create table courses_categories
(
    id        serial primary key,
    category_id int  not null,
    course_id uuid not null,
    CONSTRAINT  unique_category_course unique (category_id, course_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS uni_courses_categories ON public.courses_categories USING btree (category_id);
CREATE UNIQUE INDEX IF NOT EXISTS uni_courses_categories ON public.courses_categories USING btree (course_id);


create table courses_sub_categories
(
    id        serial primary key,
    category_id int  not null,
    course_id uuid not null,
    CONSTRAINT  unique_sub_category_course unique (category_id, course_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS uni_courses_categories ON public.courses_sub_categories USING btree (category_id);
CREATE UNIQUE INDEX IF NOT EXISTS uni_courses_categories ON public.courses_sub_categories USING btree (course_id);


create table courses_topics
(
    id        serial primary key,
    topic_id int  not null,
    course_id uuid not null,
    CONSTRAINT  unique_topic_course unique (topic_id, course_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS uni_courses_categories ON public.courses_topics USING btree (topic_id);
CREATE UNIQUE INDEX IF NOT EXISTS uni_courses_categories ON public.courses_topics USING btree (course_id);


