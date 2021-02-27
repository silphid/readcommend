# --- App build-only stage ---

FROM golang:1.16-buster as build

# For much faster builds, download go modules first because that's the longest
# part and we don't want to invalidate that work each time we change source files
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build binaries
COPY src src
RUN go build -o bin/server ./src/server
RUN go build -o bin/migrate ./src/migrations

# --- App run-time stage ---

FROM debian:buster

# Create and use non-root user for increased security
RUN mkdir /app \
  && groupadd -g 5000 -r app \
  && useradd -d /app --no-log-init -r -g app -c "Restricted app user" -u 5000 app
USER app

# Copy binaries from previous stage and migrations
WORKDIR /app
COPY --from=build /app/bin/* ./
COPY src/migrations/*.sql ./

# Env vars shared by both binaries
ENV ENV ""
ENV DB_URL ""
ENV LOG_LEVEL ""

# Env vars specific to server binary
ENV PORT ""

# Start "server" by default, but "migrate" is also available
CMD [ "./server" ]
