FROM mysql/mysql-server:latest

ENV MYSQL_DATABASE usagemeter

ADD sql-scripts/createUsers.sql /docker-entrypoint-initdb.d

