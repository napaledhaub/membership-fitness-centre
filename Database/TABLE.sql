CREATE TABLE members (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT null unique,
    email VARCHAR(100) NOT NULL unique,
    password VARCHAR(255) not null,
    isverified bool not null,
    verificationtoken VARCHAR(36),
    tokencreatedat timestamp,
    packageid int,
    expiredate timestamp,
    CONSTRAINT fk_packages
        FOREIGN KEY (packageid)
        REFERENCES packages (id)
);

CREATE TABLE packages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT null unique,
    interval int NOT NULL
);

insert into packages (name, interval)
values ('Bronze', 1);
insert into packages (name, interval)
values ('Silver', 3);
insert into packages (name, interval)
values ('Gold', 6);
