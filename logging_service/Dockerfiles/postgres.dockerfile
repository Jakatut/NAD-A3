FROM library/postgres

ADD Dockerfiles/init-db.sh /docker-entrypoint-initdb.d/

EXPOSE 5432