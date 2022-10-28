FROM postgres:12
ENV POSTGRES_DB phonedb
ENV POSTGRES_USER phone
ENV POSTGRES_HOST_AUTH_METHOD trust
CMD ["postgres", "-c", "fsync=off", "-c", "log_statement=all"]
EXPOSE 5432