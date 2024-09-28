FROM golang:1.22-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go clean --modcache && \
    go mod download && \
    go mod verify

COPY . .
RUN go build -v -o app migrations/cmd/application/main.go

FROM golang:1.22-alpine

WORKDIR /application

COPY --from=builder /build/app /application

CMD ./app --migrate="up"