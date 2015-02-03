create table texts
(
    ID int,
    name varchar(20),
    submitTime datetime,
    submitText text
);
insert into texts (ID,name,submitTime,submitText) values (1,'server',now(),"WelCome to 留言板");