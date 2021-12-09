create table channels (
    id serial primary key ,
    name VARCHAR,
    access_key VARCHAR
);

create table records (
    id serial primary key,
    channel_id int REFERENCES channels(id) not null,
    channel_one float DEFAULT 0,
    channel_two float DEFAULT 0,
    channel_three float DEFAULT 0,
    channel_four float DEFAULT 0
);