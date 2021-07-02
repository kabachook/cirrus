FROM node:16-slim AS frontend-build

WORKDIR /app

# Cache deps
COPY ./frontend/package.json ./frontend/yarn.lock ./
RUN yarn

COPY ./frontend/ ./
RUN yarn build

FROM golang:1.16 AS bin-build

WORKDIR /app

# Cache deps
COPY go.* ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -ldflags '-w -extldflags "-static"' -o bin/cirrus main.go

FROM scratch

WORKDIR /app/build
COPY --from=frontend-build /app/build/ ./

WORKDIR /app
COPY --from=bin-build /app/bin/cirrus ./

EXPOSE 3232
ENTRYPOINT [ "/app/cirrus" ]