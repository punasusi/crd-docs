# Developing

## Using Postgres Docker Image

The easiest way to get started developing locally is with the official [Postgres
Docker image](https://hub.docker.com/_/postgres).

1. Start docker container in background:

```
docker run -d --rm \
   --name dev-postgres \
   -e POSTGRES_PASSWORD=password \
   -p 5432:5432 postgres
```

2. Setup doc database and tables:

```
psql -h 127.0.0.1 -U postgres -d postgres -a -f schema/crds_up.sql
```

3. Setup doc database and tables:

```
psql -h 127.0.0.1 -U postgres -d postgres -a -f schema/crds_up.sql
```
