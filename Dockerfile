FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /bin/autoslot .

FROM alpine:3.21

WORKDIR /app
RUN adduser -D appuser && mkdir -p /data && chown -R appuser:appuser /data /app

COPY --from=build /bin/autoslot /app/autoslot
COPY templates /app/templates
COPY static /app/static
COPY data/schema.sql /app/data/schema.sql
COPY data/seed.sql /app/data/seed.sql

ENV DB_PATH=/data/autoslot.db
ENV APP_ADDR=:8080
ENV ADMIN_USER=admin
ENV ADMIN_PASS=admin123

EXPOSE 8080
USER appuser

CMD ["/app/autoslot"]
