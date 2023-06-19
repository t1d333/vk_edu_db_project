FROM golang:1.20.4-alpine3.17 as builder
COPY go.mod go.sum /forum/
WORKDIR /forum
RUN go mod download
COPY . /forum

RUN go build -o app ./cmd/forum/main.go


FROM postgres:14 as db
USER postgres
RUN pg_createcluster 14 main && \
    /etc/init.d/postgresql start &&\
    psql --command "CREATE USER forum WITH SUPERUSER PASSWORD 'password';" &&\
    createdb -O forum forum && \
    /etc/init.d/postgresql stop 


ENV POSTGRES_USER forum
ENV POSTGRES_DB forum
ENV POSTGRES_PASSWORD password
ENV POSTGRES_HOST localhost
ENV POSTGRES_PORT 5432

USER root
COPY --from=builder /forum/app app
COPY ./scripts/init.sql .
ENV PGPASSWORD password 
CMD service postgresql start && psql -h localhost -d forum -U forum -p 5432 -a -q -f ./init.sql && ./app








