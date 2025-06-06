FROM golang:1.24-alpine3.20

WORKDIR /

RUN ls -la

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

CMD ["go", "run", "/cmd/worker/single_replica/main.go"]