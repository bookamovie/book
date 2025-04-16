FROM golang:1.23 AS build

WORKDIR /book

COPY . .

RUN go mod download
RUN go build -o ./build/book ./cmd/book/main.go && \
    go build -o ./build/migrator ./cmd/migrator/main.go

FROM debian:bookworm-slim

LABEL org.opencontainers.image.source https://github.com/bookamovie/book
LABEL org.opencontainers.image.description Book Microservice Image
LABEL org.opencontainers.image.licenses MIT

WORKDIR /book

COPY --from=build book/build/book .
COPY --from=build book/build/migrator tools/

COPY --from=build book/deployments/docker/config config/
COPY --from=build book/migrations migrations/
COPY --from=build book/storage storage/

EXPOSE 5092

ENTRYPOINT [ "./book" ]
