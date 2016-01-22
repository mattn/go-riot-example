drop table todo;
create table todo (id serial, title text, done bool default false, createat timestamp default 'now');
