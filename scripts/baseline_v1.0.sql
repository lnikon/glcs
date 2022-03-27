create table "Computations"
(
    UID serial not null,
    Name varchar(128) not null,
    Algorithm varchar(128) not null,
    VertexCount varchar(128) not null,
    Density varchar(128),
    Replicas integer,
    StartTime timestamp not null,
    EndTime timestamp,
    Status varchar(128),
    Result text,
    constraint Computations_PrimaryKey primary key (UID)
);

create unique index Computations_Name_UniqueIndex on "public"."Computations"(Name);

