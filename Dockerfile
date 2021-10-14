# build stage
FROM golang:alpine AS build-env
RUN apk add --update gcc make musl-dev
WORKDIR /src/
COPY . ./
RUN make project-utils
RUN make all

# final stage
FROM alpine
RUN apk add make
WORKDIR /bot
COPY Makefile ./
COPY migrations ./migrations/
COPY data/sqlite/bot.sqlite3 ./data/sqlite/bot.sqlite3
COPY --from=build-env /src/bin/* ./
COPY --from=build-env /go/bin/migrate /bin/
RUN mkdir -p data/sqlite
ENTRYPOINT make migrate-up && ls data/sqlite && ./bot -t $BOT_TOKEN -db data/sqlite/bot.sqlite3

