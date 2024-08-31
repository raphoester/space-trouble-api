CREATE TABLE bookings (
    ID            text not null primary key,
    FirstName     text not null,
    LastName      text not null,
    Gender        text not null,
    Birthday      text not null,
    DestinationID text not null,
    LaunchpadID   text not null,
    LaunchDate    text not null
);