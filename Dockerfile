FROM golang:1.20.4-alpine3.17 as builder
COPY go.mod go.sum /forum/
WORKDIR /forum
RUN \
    --mount=type=cache,target=/go/pkg/mod/ \
    go mod download
COPY . /forum

RUN \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o app ./cmd/forum/main.go


FROM postgres:14 as db
USER postgres
COPY ./scripts/init.sql .
RUN initdb && \
    pg_ctl start && \
    createuser -s user && \
    createdb -O user forum && \
    psql -U user -f './init.sql' -d forum \
    pg_ctl stop

ENV POSTGRES_USER user
ENV POSTGRES_DB forum
ENV POSTGRES_PASSWORD password
ENV POSTGRES_HOST localhost
ENV POSTGRES_PORT 5432

COPY --from=builder /forum/app app

CMD pg_ctl start && ./app








