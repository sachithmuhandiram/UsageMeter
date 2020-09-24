FROM mysql/mysql-server:latest

COPY ./sql-scripts/ /docker-entrypoint-initdb.d/

ENV MYSQL_DATABASE usagemeter