create table users (
                       id varchar primary key,
                       username varchar not null,
                       email varchar not null unique,
                       password varchar not null,
                       created_at timestamp,
                       updated_at timestamp
);

create table photos (
                        id varchar primary key,
                        title varchar,
                        caption varchar,
                        photo_url varchar,
                        user_id varchar unique,
                        foreign key (user_id) references users(id) on delete cascade
);
