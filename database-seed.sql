-- Create users table
DROP TABLE if exists users;
CREATE TABLE IF NOT EXISTS public.users
(
    id                serial primary key,
    user_id           uuid                     not null,
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
CREATE UNIQUE INDEX IF NOT EXISTS uni_users_user_id ON public.users USING btree (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS uni_users_email ON public.users USING btree (email);
CREATE UNIQUE INDEX IF NOT EXISTS uni_users_username ON public.users USING btree (username);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON public.users USING btree (deleted_at);


DROP TABLE if exists account_status;
CREATE TABLE IF NOT EXISTS public.account_status
(
    id               serial primary key,
    account_id       uuid not null references users (user_id) on DELETE CASCADE,
    status           varchar(128),
    note             varchar(5000),
    created_at       timestamp with time zone default current_timestamp,
    updated_at       timestamp with time zone default current_timestamp,
    is_status_active boolean                  default true,
    CONSTRAINT unique_account_status_active unique (account_id, status, is_status_active)
);

-- Create indexes for account_status table
CREATE INDEX IF NOT EXISTS uni_account_status_account_id ON public.account_status USING btree (account_id);


DROP TABLE if exists profiles;
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
    user_id    uuid               not null,
    foreign key (user_id) references public.users (user_id)
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
    course_id    uuid unique        not null,
    created_at   timestamp with time zone default current_timestamp,
    updated_at   timestamp with time zone default current_timestamp,
    deleted_at   timestamp with time zone,
    thumbnail    text,
    title        text               not null unique,
    slug         text               not null,
    description  varchar,
    publish_date timestamp with time zone,
    price        double precision   not null
);

-- Create indexes for courses table
CREATE UNIQUE INDEX IF NOT EXISTS uni_courses_slug ON public.courses USING btree (slug);
CREATE INDEX IF NOT EXISTS idx_courses_deleted_at ON public.courses USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_courses_course_id ON public.courses USING btree (course_id);



DROP TABLE IF EXISTS authors_courses;
create table authors_courses
(
    id        serial primary key,
    author_id uuid  not null references users(user_id),
    course_id uuid not null references courses (course_id) on DELETE CASCADE,
    CONSTRAINT unique_author_course unique (author_id, course_id)
);
CREATE INDEX IF NOT EXISTS uni_authors_courses_author_id ON public.authors_courses USING btree (author_id);
CREATE INDEX IF NOT EXISTS uni_authors_courses_course_id ON public.authors_courses USING btree (course_id);


create table courses_categories
(
    id          serial primary key,
    category_id int  not null references categories (id),
    course_id   uuid not null references courses (course_id),
    CONSTRAINT unique_category_course unique (category_id, course_id)
);
CREATE INDEX IF NOT EXISTS uni_courses_categories ON public.courses_categories USING btree (category_id);
CREATE INDEX IF NOT EXISTS uni_courses_categories ON public.courses_categories USING btree (course_id);


create table courses_sub_categories
(
    id          serial primary key,
    category_id int  not null references categories (id),
    course_id   uuid not null references courses (course_id),
    CONSTRAINT unique_sub_category_course unique (category_id, course_id)
);
CREATE INDEX IF NOT EXISTS uni_courses_sub_categories ON public.courses_sub_categories USING btree (category_id);
CREATE INDEX IF NOT EXISTS uni_courses_sub_categories ON public.courses_sub_categories USING btree (course_id);


create table courses_topics
(
    id        serial primary key,
    topic_id  int  not null references categories (id),
    course_id uuid not null references courses (course_id),
    CONSTRAINT unique_topic_course unique (topic_id, course_id)
);
CREATE INDEX IF NOT EXISTS uni_courses_topics ON public.courses_topics USING btree (topic_id);
CREATE INDEX IF NOT EXISTS uni_courses_topics ON public.courses_topics USING btree (course_id);


DROP table if exists roles;
create table roles
(
    id          serial primary key,
    role_id     uuid unique  not null,
    name        varchar(512) not null,
    description varchar(2000),
    status      varchar(8)   not null    default 'active',
    slug        varchar(512) not null,
    created_at  timestamp with time zone default current_timestamp,
    updated_at  timestamp with time zone default current_timestamp,
    deleted_at  timestamp with time zone
);

CREATE UNIQUE INDEX IF NOT EXISTS uni_roles_title ON public.roles USING btree (name);
CREATE UNIQUE INDEX IF NOT EXISTS uni_roles_slug ON public.roles USING btree (slug);
CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON public.roles USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_roles_status ON public.roles USING btree (status);



DROP table if exists users_roles;
create table users_roles
(
    id      serial primary key,
    user_id uuid not null references users (user_id),
    role_id uuid not null references roles (role_id),
    CONSTRAINT unique_user_role unique (user_id, role_id)
);

CREATE INDEX IF NOT EXISTS uni_users_roles_user_id ON public.users_roles USING btree (user_id);
CREATE INDEX IF NOT EXISTS uni_users_roles_role_id ON public.users_roles USING btree (role_id);


/*Rating & Review */
/*rate: number,
title: string,
summary: string,
images: []
username: string,
createdAt: string
avatar: string
*/

DROP TABLE IF EXISTS reviews;
create table reviews
(
    id         serial primary key,
    title      varchar(512),
    summary    varchar(10000),
    rate       int8,
    user_id    uuid not null references users (user_id),
    course_id  uuid not null references courses (course_id),
    created_at timestamptz default current_timestamp,
    deleted_at timestamptz default null
);

CREATE INDEX IF NOT EXISTS i_reviews_user_id ON public.reviews USING btree (user_id);
CREATE INDEX IF NOT EXISTS i_reviews_course_id ON public.reviews USING btree (course_id);
CREATE INDEX IF NOT EXISTS i_reviews_deleted_at ON public.reviews USING btree (deleted_at);



select courses.*,
       (select title from courses_categories c where courses.course_id = c.course_id)       as categories,
       (select title from courses_sub_categories sc where courses.course_id = sc.course_id) as sub_categories,
       (select title from courses_topics ct where courses.course_id = ct.course_id)         as topics
from courses
         join authors_courses ac on courses.course_id = ac.course_id
where ac.author_id = 6;


-- get users roles
select u.email, u.id as user_id, json_agg(DISTINCT r.role_id) AS "roles"
from users u
         left join public.users_roles ur
                   on ur.user_id = u.id
         left join roles r on ur.role_id = r.role_id
group by u.email, u.id;



select users.email,
       users.user_id,

       (select jsonb_agg(jsonb_build_object(
               'status', ass.status,
               'is_status_active', ass.is_status_active)
               ) AS account_status
        from account_status ass
        where ass.account_id = users.user_id
          AND is_status_active = true)

from users
         join users_roles ur
              on users.user_id = ur.user_id
         join public.roles r
              on ur.role_id = r.role_id
where r.slug = 'instructor';



SELECT courses.id                               as id,
       ac.course_id                             as course_id,
       title,
       slug,
       description,
       thumbnail,
       price,
       created_at,

       (select jsonb_agg(DISTINCT cs.category_id)
        from courses_categories cs
        where courses.course_id = cs.course_id) as categories,
       (select jsonb_agg(DISTINCT sc.category_id)
        from courses_sub_categories sc
        where courses.course_id = sc.course_id) as sub_categories,
       (select jsonb_agg(DISTINCT ct.topic_id)
        from courses_topics ct
        where courses.course_id = ct.course_id) as topics,
       (select jsonb_agg(DISTINCT ac.author_id)
        from authors_courses ac
        where ac.course_id = courses.course_id) as authors

FROM courses
         join authors_courses ac on courses.course_id = ac.course_id
where courses.slug = 'quaerat-quidem-non-q'

