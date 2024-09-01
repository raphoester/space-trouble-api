CREATE TABLE bookings (
    id            text not null primary key,
    first_name     text not null,
    last_name      text not null,
    gender        text not null,
    birthday      text not null,
    destination_id text not null,
    launchpad_id   text not null,
    launch_date    text not null
);