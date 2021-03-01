# --- Build go app ---

FROM golang:1.16-buster as build-go

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

# --- Build node app ---

FROM node:14 as build-node

# For much faster builds, install node packages first because that's the longest
# part and we don't want to invalidate that work each time we change source files
WORKDIR /app
COPY app/package.json .
RUN npm install

# Build node app
COPY app/ .
RUN npm run build

# --- Final app run-time ---

FROM debian:buster

# Create and use non-root user for increased security
RUN mkdir /app \
  && groupadd -g 5000 -r app \
  && useradd -d /app --no-log-init -r -g app -c "Restricted app user" -u 5000 app
USER app

# Copy outputs previous stages and migrations
WORKDIR /app
COPY --from=build-go /app/bin/* ./
COPY --from=build-node /app/dist/* ./app/dist/
COPY src/migrations/*.sql ./
COPY dist/* ./

ENV ENV ""
ENV DB_URL "postgres://postgres:password123@db:5432/readcommend?sslmode=disable"
ENV DB_HOST_PORT "db:5432"
ENV LOG_LEVEL "INFO"
ENV PORT "5000"

CMD [ "./start.sh" ]
